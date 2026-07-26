import React, {useCallback, useEffect, useRef, useState} from "react";
import {App, Card} from "antd";
import {Body} from "sinking-antd";
import {type MindRef, type MindSaveReason, type MindStats, type MindValue,} from "@/pages/components/mind";
import {restoreMindVersion, saveMindData} from "@/service/user/mind";
import {VERSION_AUTO, VERSION_MANUAL} from "./constants";
import useStyles from "./styles";
import type {MindDocument, MindSaveTask, MindVersionRecord,} from "./types";
import {
    cloneValue,
    createBlankValue,
    downloadDataUrl,
    getMindValueSignature,
    getResponseData,
    toMindDocument,
    toMindVersion,
} from "./utils";
import Category from "./components/category";
import CategoryForm from "./components/category-form";
import Data from "./components/data";
import Editor from "./components/editor";
import Header from "./components/header";
import Info from "./components/info";
import Preview from "./components/preview";
import Share from "./components/share";
import ShareForm from "./components/share-form";
import Version from "./components/version";
import {useMindCategories, useMindDocuments, useMindPageHeight, useMindShares, useMindVersions,} from "./hooks";

const MindPage: React.FC = () => {
    const {styles} = useStyles();
    const {message, modal} = App.useApp();
    const editorRef = useRef<MindRef | null>(null);
    const selectedDocumentRef = useRef<MindDocument | undefined>(undefined);
    const editorValueRef = useRef<MindValue>(createBlankValue("中心主题"));
    const editorDirtyRef = useRef(false);
    const lastSavedRef = useRef(new Map<string, string>());
    const saveTasksRef = useRef<MindSaveTask[]>([]);
    const saveProcessingRef = useRef(false);
    const activeSaveTaskRef = useRef<MindSaveTask | undefined>(undefined);
    const saveDrainWaitersRef = useRef<Array<{dataId: string; resolve: () => void}>>([]);
    const processSaveTasksRef = useRef<() => void>(() => undefined);
    const [selectedDocument, setSelectedDocument] = useState<MindDocument>();
    const [editorValue, setEditorValue] = useState<MindValue>(() => createBlankValue("中心主题"));
    const [editorRevision, setEditorRevision] = useState(0);
    const [editorClosing, setEditorClosing] = useState(false);
    const versions = useMindVersions({selectedDocumentRef});
    const {pageRef, pageStyle} = useMindPageHeight(selectedDocument?.id);

    useEffect(() => {
        selectedDocumentRef.current = selectedDocument;
    }, [selectedDocument]);

    useEffect(() => {
        editorValueRef.current = editorValue;
    }, [editorValue]);

    const showEditor = useCallback((record: MindDocument, value: MindValue) => {
        const nextValue = cloneValue(value);
        const nextRecord = {...record, value: nextValue};
        editorValueRef.current = nextValue;
        editorDirtyRef.current = false;
        lastSavedRef.current.clear();
        lastSavedRef.current.set(record.id, getMindValueSignature(nextValue));
        setEditorValue(nextValue);
        setEditorClosing(false);
        setSelectedDocument(nextRecord);
        versions.reset();
        setEditorRevision((revision) => revision + 1);
    }, [versions.reset]);

    const documentActionsRef = useRef<{
        moveCategoryToAll: (categoryId: string) => void;
        reloadDocuments: () => void;
    }>({
        moveCategoryToAll: () => undefined,
        reloadDocuments: () => undefined,
    });
    const handleCategoryDeleted = useCallback((categoryId: string) => {
        documentActionsRef.current.moveCategoryToAll(categoryId);
        documentActionsRef.current.reloadDocuments();
    }, []);
    const categories = useMindCategories({onDeleted: handleCategoryDeleted});
    const handleDocumentUpdated = useCallback((record: MindDocument) => {
        setSelectedDocument((current) => current?.id === record.id ? {...current, ...record} : current);
    }, []);
    const handleDocumentDeleted = useCallback((record: MindDocument) => {
        if (selectedDocumentRef.current?.id !== record.id) return;
        setSelectedDocument(undefined);
        versions.reset();
    }, [versions.reset]);
    const documents = useMindDocuments({
        activeCategoryId: categories.activeCategoryId,
        onOpen: showEditor,
        onUpdated: handleDocumentUpdated,
        onDeleted: handleDocumentDeleted,
    });
    useEffect(() => {
        documentActionsRef.current = {
            moveCategoryToAll: documents.moveCategoryToAll,
            reloadDocuments: documents.reloadDocuments,
        };
    }, [documents.moveCategoryToAll, documents.reloadDocuments]);
    const shares = useMindShares({onChanged: documents.reloadDocuments});

    const executeSaveTask = useCallback(async (task: MindSaveTask) => {
        const {dataId, value, signature, saveType, reason} = task;
        if (saveType === VERSION_AUTO && lastSavedRef.current.get(dataId) === signature) {
            if (selectedDocumentRef.current?.id === dataId
                && getMindValueSignature(editorValueRef.current) === signature) {
                editorDirtyRef.current = false;
            }
            return;
        }
        const response = await saveMindData({body: {
            id: Number(dataId),
            value: JSON.stringify(value),
            type: saveType,
            remark: reason === "shortcut" ? "快捷键保存" : "",
        }});
        const savedVersion = toMindVersion(getResponseData<any>(response, "保存导图失败"));
        lastSavedRef.current.set(dataId, signature);
        if (selectedDocumentRef.current?.id === dataId
            && getMindValueSignature(editorValueRef.current) === signature) {
            editorDirtyRef.current = false;
        }
        const updateTime = savedVersion.updateTime || savedVersion.createTime;
        documents.updateDocumentTime(dataId, updateTime);
        setSelectedDocument((current) => current?.id === dataId ? {...current, updateTime} : current);
        versions.recordSavedVersion(savedVersion, saveType);
        if (reason !== "auto") {
            message.success("导图已保存");
        }
    }, [documents.updateDocumentTime, message, versions.recordSavedVersion]);

    const hasPendingSave = useCallback((dataId: string) => {
        return activeSaveTaskRef.current?.dataId === dataId
            || saveTasksRef.current.some((task) => task.dataId === dataId);
    }, []);

    const resolveSaveDrainWaiters = useCallback(() => {
        const pending: typeof saveDrainWaitersRef.current = [];
        saveDrainWaitersRef.current.forEach((waiter) => {
            if (hasPendingSave(waiter.dataId)) {
                pending.push(waiter);
            } else {
                waiter.resolve();
            }
        });
        saveDrainWaitersRef.current = pending;
    }, [hasPendingSave]);

    const processSaveTasks = useCallback(() => {
        if (saveProcessingRef.current) return;
        const task = saveTasksRef.current.shift();
        if (!task) {
            resolveSaveDrainWaiters();
            return;
        }
        saveProcessingRef.current = true;
        activeSaveTaskRef.current = task;
        void executeSaveTask(task).then(() => {
            task.waiters.forEach(({resolve}) => resolve());
        }).catch((error) => {
            task.waiters.forEach(({reject}) => reject(error));
        }).finally(() => {
            activeSaveTaskRef.current = undefined;
            saveProcessingRef.current = false;
            processSaveTasksRef.current();
            resolveSaveDrainWaiters();
        });
    }, [executeSaveTask, resolveSaveDrainWaiters]);

    useEffect(() => {
        processSaveTasksRef.current = processSaveTasks;
        processSaveTasks();
    }, [processSaveTasks]);

    const persistDocument = useCallback((value?: MindValue, stats?: MindStats, reason: MindSaveReason = "manual") => {
        const document = selectedDocumentRef.current;
        if (!document || (reason === "auto" && !editorDirtyRef.current)) return Promise.resolve();
        const dataId = document.id;
        const nextValue = cloneValue(value || editorRef.current?.getValue() || editorValueRef.current);
        const signature = getMindValueSignature(nextValue);
        const saveType = reason === "auto" ? VERSION_AUTO : VERSION_MANUAL;
        void stats;
        return new Promise<void>((resolve, reject) => {
            const tasks = saveTasksRef.current;
            const pending = tasks[tasks.length - 1];
            if (saveType === VERSION_AUTO && pending?.saveType === VERSION_AUTO && pending.dataId === dataId) {
                pending.value = nextValue;
                pending.signature = signature;
                pending.reason = reason;
                pending.waiters.push({resolve, reject});
            } else {
                tasks.push({
                    dataId,
                    value: nextValue,
                    signature,
                    saveType,
                    reason,
                    waiters: [{resolve, reject}],
                });
            }
            processSaveTasksRef.current();
        });
    }, []);

    const waitForDocumentSaves = useCallback((dataId: string) => {
        if (!hasPendingSave(dataId)) return Promise.resolve();
        return new Promise<void>((resolve) => {
            saveDrainWaitersRef.current.push({dataId, resolve});
        });
    }, [hasPendingSave]);

    const saveEditor = () => {
        const task = editorRef.current?.save("manual");
        if (task) {
            void task.catch((error) => message.error(error instanceof Error ? error.message : "导图保存失败"));
        }
    };

    const restoreVersion = async (record: MindVersionRecord) => {
        const document = selectedDocumentRef.current;
        if (!document) return;
        setEditorClosing(true);
        try {
            if (editorDirtyRef.current) {
                await editorRef.current?.save("auto");
            }
            await waitForDocumentSaves(document.id);
            const response = await restoreMindVersion({body: {
                data_id: Number(document.id),
                id: Number(record.id),
            }});
            const restored = toMindDocument(getResponseData<any>(response, "恢复版本失败"));
            if (!restored.value) throw new Error("恢复后的导图内容为空");
            showEditor(restored, restored.value);
            documents.reloadDocuments();
            message.success("版本已恢复");
        } finally {
            setEditorClosing(false);
        }
    };

    const confirmRestoreVersion = (record: MindVersionRecord) => {
        if (record.value
            && getMindValueSignature(editorRef.current?.getValue() || editorValueRef.current)
            === getMindValueSignature(record.value)) {
            message.info("当前内容已与该版本一致");
            return;
        }
        modal.confirm({
            title: "恢复此版本？",
            content: `将使用 ${record.createTime} 保存的内容替换当前画布。`,
            okText: "确认",
            cancelText: "取消",
            onOk: async () => {
                try {
                    await restoreVersion(record);
                } catch (error) {
                    message.error(error instanceof Error ? error.message : "版本恢复失败");
                }
            },
        });
    };

    const arrangeEditor = () => {
        editorRef.current?.layout();
    };

    const exportEditorImage = async (type: "png" | "jpg") => {
        const editor = editorRef.current;
        if (!editor || !selectedDocument) {
            return;
        }
        const dataUrl = type === "png" ? await editor.exportPng() : await editor.exportJpg();
        downloadDataUrl(dataUrl, `${selectedDocument.title}.${type}`);
        message.success(`${type.toUpperCase()} 已导出`);
    };

    const copyEditorImage = () => {
        const task = editorRef.current?.copyImage();
        if (task) {
            void task.catch((error) => message.error(error instanceof Error ? error.message : "复制图片失败"));
        }
    };

    const categorySelectOptions = categories.getCategorySelectOptions(documents.editingDocument?.categoryId);
    const closeEditor = async () => {
        if (editorClosing) return;
        const document = selectedDocumentRef.current;
        if (!document) {
            documents.cancelOpenDocument();
            return;
        }
        setEditorClosing(true);
        try {
            if (editorDirtyRef.current) {
                await editorRef.current?.save("auto");
            }
            await waitForDocumentSaves(document.id);
            setSelectedDocument(undefined);
            versions.reset();
        } catch (error) {
            message.error(error instanceof Error ? error.message : "导图保存失败");
        } finally {
            setEditorClosing(false);
        }
    };
    return (
        <Body space={false}>
            <div ref={pageRef} className={styles.page} style={pageStyle}>
                <div className={styles.library}>
                    <Card className={styles.libraryCard} variant="borderless">
                        <Category
                            styles={styles}
                            categories={categories.categories}
                            activeCategoryId={categories.activeCategoryId}
                            loading={categories.categoryLoading}
                            loadFailed={categories.categoryLoadFailed}
                            onSelect={categories.setActiveCategoryId}
                            onLoadMore={() => void categories.loadMoreCategories()}
                            onCreate={categories.openCreateCategory}
                            onEdit={categories.openEditCategory}
                            onDelete={categories.deleteCategory}
                            onRetry={categories.retryCategories}
                        />

                        <div className={styles.content}>
                            <Header
                                styles={styles}
                                keyword={documents.keywordInput}
                                sortType={documents.sortType}
                                onKeywordChange={documents.setKeywordInput}
                                onSearch={documents.searchDocuments}
                                onSortChange={documents.setSortType}
                                onCreate={documents.openCreate}
                            />
                            <Data
                                styles={styles}
                                documents={documents.documents}
                                loading={documents.documentsLoading}
                                total={documents.documentTotal}
                                page={documents.page}
                                pageSize={documents.pageSize}
                                openingDocumentId={documents.documentOpeningId}
                                onOpen={documents.openDocument}
                                onEdit={documents.openEditDocument}
                                onCopy={documents.duplicateDocument}
                                onShare={shares.openShareManager}
                                onDelete={documents.deleteDocument}
                                onPaginationChange={documents.changePagination}
                            />
                        </div>
                    </Card>
                </div>
            </div>

            <Editor
                styles={styles}
                open={Boolean(selectedDocument || documents.documentOpeningId)}
                loading={Boolean(documents.documentOpeningId && !selectedDocument)}
                document={selectedDocument}
                editorRef={editorRef}
                value={editorValue}
                revision={editorRevision}
                closing={editorClosing}
                onExportPng={() => exportEditorImage("png").catch((error) => {
                    message.error(error instanceof Error ? error.message : "PNG 导出失败");
                })}
                onExportJpg={() => exportEditorImage("jpg").catch((error) => {
                    message.error(error instanceof Error ? error.message : "JPG 导出失败");
                })}
                onCopyImage={copyEditorImage}
                onClose={closeEditor}
                onSave={saveEditor}
                onArrange={arrangeEditor}
                onOpenVersionHistory={versions.openHistory}
                onOpenShare={shares.openShareManager}
                onDataChange={(value) => {
                    editorValueRef.current = value;
                    editorDirtyRef.current = true;
                    setEditorValue(value);
                }}
                onDirty={() => {
                    editorDirtyRef.current = true;
                }}
                onPersist={persistDocument}
            />

            <Version
                styles={styles}
                open={versions.open}
                document={selectedDocument}
                filter={versions.filter}
                items={versions.items}
                loading={versions.loading}
                loadFailed={versions.loadFailed}
                onClose={versions.close}
                onLoadMore={versions.loadMore}
                onFilterChange={versions.changeFilter}
                onPreview={versions.openPreview}
            />

            <Share
                styles={styles}
                open={shares.shareOpen}
                document={shares.shareDocument}
                items={shares.shareItems}
                loading={shares.shareLoading}
                loadFailed={shares.shareLoadFailed}
                onClose={shares.closeShareManager}
                onCreate={shares.openCreateShare}
                onLoadMore={shares.loadShares}
                onCopyLink={shares.copyShareLink}
                onEdit={shares.openEditShare}
                onRemove={shares.removeShare}
            />

            <Preview
                styles={styles}
                open={versions.previewOpen}
                loading={versions.previewLoading}
                version={versions.preview}
                onRestore={confirmRestoreVersion}
                onClose={versions.closePreview}
            />

            <Info
                styles={styles}
                form={documents.form}
                open={documents.createOpen}
                editingDocument={documents.editingDocument}
                submitting={documents.formSubmitting}
                categoryOptions={categorySelectOptions}
                categoryLoading={categories.categoryLoading}
                onCategoryPopupScroll={categories.handleCategoryPopupScroll}
                onSave={documents.saveDocumentInfo}
                onCancel={documents.closeDocumentForm}
            />

            <ShareForm
                styles={styles}
                form={shares.shareForm}
                open={shares.shareFormOpen}
                editingShare={shares.editingShare}
                submitting={shares.shareSubmitting}
                onSave={shares.saveShare}
                onCancel={shares.closeShareForm}
            />

            <CategoryForm
                styles={styles}
                form={categories.categoryForm}
                open={categories.categoryOpen}
                category={categories.categoryRecord}
                submitting={categories.categorySubmitting}
                onSave={categories.saveCategory}
                onCancel={categories.closeCategoryForm}
            />
        </Body>
    );
};

export default MindPage;
