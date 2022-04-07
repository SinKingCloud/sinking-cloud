import React, {useEffect, useRef, useState} from 'react';
import {Button, Dropdown, Form, Menu, message, Modal, Select, Table, Tooltip, Upload} from 'antd';
import {FooterToolbar, PageContainer} from '@ant-design/pro-layout';
import ProTable from '@ant-design/pro-table';
import {getPageControls, getParams} from "@/util/page";
import {DownOutlined, ExclamationCircleOutlined, UploadOutlined} from "@ant-design/icons";
import {deleteFile, getFileList} from "@/services/system/file";
import {getUploadUrl} from "@/services/common/upload";

export default (): React.ReactNode => {
  /**
   * 控件权限
   */
  const [controls, setControls] = useState({});
  useEffect(() => {
    getPageControls().then((d) => {
      setControls(d)
    })
  }, []);
  const [isModalUploadVisible, setIsModalUploadVisible] = useState(false);
  const [isUploadLoading, setIsUploadLoading] = useState(false);
  const [uploadType, setUploadType] = useState(1);
  /**
   * 表单处理
   */
  const actionRef = useRef();
  const ref = useRef();
  const [selectedRowsState, setSelectedRows] = useState([]);
  /**
   * 获取选中项
   */
  const getSelectedIds = () => {
    let ids: any[] = [];
    selectedRowsState.map((k: any) => {
      ids.push(k.id);
    });
    return ids;
  };

  /**
   * 删除数据
   * @param ids id列表
   */
  const deleteSubmit = (ids: any[]) => {
    Modal.confirm({
      title: '您确定要删除该数据吗?',
      icon: <ExclamationCircleOutlined/>,
      content: '删除后该数据不可恢复',
      okType: 'danger',
      onOk() {
        deleteFile({ids: ids}).then((r) => {
          if (r.code != 200) {
            message.error(r.message || "请求失败").then()
          } else {
            // @ts-ignore
            actionRef.current.reloadAndRest()
            message.success(r.message).then()
          }
        })
      },
    });
  }

  /**
   * 批量删除
   */
  const modalDelete = () => {
    deleteSubmit(getSelectedIds())
  }

  const copy = (text: string) => {
    var oInput = document.createElement('input');
    oInput.value = text;
    document.body.appendChild(oInput);
    oInput.select();
    document.execCommand("Copy");
    oInput.className = 'oInput';
    oInput.style.display = 'none';
    // @ts-ignore
    oInput.parentNode.removeChild(oInput);
  }

  /**
   * table表格渲染
   */
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      tip: '记录ID',
      hideInSearch: true,
      sorter: true,
    },
    {
      title: '文件类别',
      dataIndex: 'type',
      tip: '文件所属类别',
      valueEnum: {
        0: {
          text: '用户头像',
          color: "green",
        },
        1: {
          text: '临时文件',
          color: "yellow",
        },
      },
    },
    {
      title: '文件MD5',
      dataIndex: 'hash',
      tip: '文件指纹MD5值',
    },
    {
      title: '文件类型',
      dataIndex: 'ext',
      tip: '文件类型',
      valueEnum: {
        ".png": {
          text: 'png',
        },
        ".jpg": {
          text: 'jpg',
        },
        ".gif": {
          text: 'gif',
        },
        ".excel": {
          text: 'excel',
        },
        ".zip": {
          text: 'zip',
        },
        ".rar": {
          text: 'rar',
        },
        ".gz": {
          text: 'gz',
        },
        ".m3u8": {
          text: 'm3u8',
        },
        ".mp3": {
          text: 'mp3',
        },
        ".mp4": {
          text: 'mp4',
        },
      },
    },
    {
      title: '文件名',
      dataIndex: 'name',
      tip: '文件名称',
      hideInSearch: true,
      copyable: true,
      render: (text: string, record: any) => {
        return (
          <Tooltip title={record?.name}>
            <span style={{fontSize: "small", cursor: "pointer"}} onClick={() => {
              copy(record?.name);
              message.success("已复制到剪贴板")
            }}>点击复制</span>
          </Tooltip>);
      }
    },
    {
      title: '外链地址',
      dataIndex: 'url',
      tip: '外链地址',
      hideInSearch: true,
      render: (text: string, record: any) => {
        return (
          <Tooltip title={record?.url}>
            <span style={{fontSize: "small", cursor: "pointer"}} onClick={() => {
              copy(record?.url);
              message.success("已复制到剪贴板")
            }}>点击复制</span>
          </Tooltip>);
      }
    },
    {
      title: '文件大小',
      dataIndex: 'size',
      tip: '文件大小',
      hideInSearch: true,
      render: (text: string, record: any) => {
        let kb = record?.size / 1024
        if (kb < 1024) {
          return kb.toFixed(2) + "KB";
        }
        return (kb / 1024).toFixed(2) + "MB";
      }
    },
    {
      title: '修改时间',
      valueType: 'dateTime',
      dataIndex: 'update_time',
      tip: '上次修改时间',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: '上传时间',
      valueType: 'dateTime',
      dataIndex: 'create_time',
      tip: '文件上传时间',
      sorter: true,
      hideInSearch: true,
    },
    {
      title: '修改时间',
      valueType: 'dateTimeRange',
      dataIndex: 'update_time',
      tip: '上次修改时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            update_time_start: value[0],
            update_time_end: value[1],
          };
        },
      },
    },
    {
      title: '上传时间',
      valueType: 'dateTimeRange',
      dataIndex: 'create_time',
      tip: '文件上传时间',
      hideInTable: true,
      search: {
        transform: (value: any) => {
          return {
            create_time_start: value[0],
            create_time_end: value[1],
          };
        },
      },
    },
    {
      title: '操作',
      dataIndex: 'Option',
      hideInSearch: true,
      editable: false,
      render: (text: string, record: any) => {
        return <Dropdown overlay={(
          <Menu>
            <Menu.Item key={"download"}>
              <a style={{fontSize: "small"}} onClick={() => {
                window.open(record?.url);
              }}>
                下 载
              </a>
            </Menu.Item>
            <Menu.Item key={"delete"} hidden={controls["system.file:delete"] == undefined}>
              <a style={{fontSize: "small"}} onClick={() => {
                // @ts-ignore
                deleteSubmit([record.id])
              }}>
                删 除
              </a>
            </Menu.Item>
          </Menu>
        )} trigger={['click']} placement="bottomCenter" arrow={true}>
          <Button size={"small"} style={{fontSize: "small"}} onClick={e => e.preventDefault()}>操
            作 <DownOutlined/></Button>
        </Dropdown>;
      }
    },
  ];

  return (
    <PageContainer title={false}>
      <Modal key={"upload"} width={350} destroyOnClose={true} forceRender={true} title="上传文件"
             visible={isModalUploadVisible} onOk={() => {
        setIsModalUploadVisible(false);
      }} okText={"上传"} onCancel={() => {
        setIsModalUploadVisible(false);
      }} footer={[<Upload key={"upload"} name="file" data={{"type": uploadType}} showUploadList={false} action={getUploadUrl()}
                          beforeUpload={(file: any) => {
                            return true;
                          }}
                          onChange={(info: any) => {
                            if (info.file.status == "uploading") {
                              setIsUploadLoading(true);
                            } else if (info.file.status === 'done') {
                              setIsUploadLoading(false);
                              setIsModalUploadVisible(false)
                              const url = info.file.response.data.url;
                              if (url != undefined && url != "") {
                                // @ts-ignore
                                actionRef?.current.reload();
                                copy(url);
                                message.success("文件上传成功,外链已复制到剪贴板");
                              } else {
                                message.error(`文件上传失败`);
                              }
                            } else if (info.file.status === 'error') {
                              setIsUploadLoading(false);
                              message.error(`文件上传失败`);
                            }
                          }}>
        <Button type={"primary"} loading={isUploadLoading}>
          <UploadOutlined/>上传文件
        </Button>
      </Upload>]}>
        <Form name="control-hooks" key="upload_form" labelAlign="right" labelCol={{span: 6}}
              wrapperCol={{span: 18}}>
          <Form.Item name={"type"} label="文件所属" initialValue={uploadType}>
            <Select placeholder="请选择文件所属类型" onChange={(e) => {
              setUploadType(e);
            }}>
              <Select.Option value={1}>临时文件</Select.Option>
              <Select.Option value={0}>用户头像</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
      <ProTable
        // @ts-ignore
        columns={columns}
        defaultSize={"small"}
        form={{layout: "vertical"}}
        headerTitle={'接口列表'}
        actionRef={actionRef}
        formRef={ref}
        rowKey={'id'}
        options={{
          density: true,
          fullScreen: true,
          setting: true,
        }}

        rowSelection={{
          selections: [Table.SELECTION_ALL, Table.SELECTION_INVERT, Table.SELECTION_NONE],
          onChange: (_, selectedRows: any) => {
            setSelectedRows(selectedRows);
          },
        }}
        request={async (params, sort) => {
          const data = await getFileList(getParams(params, sort));
          if (data.code != 200) {
            message.error(data.message);
          }
          return {
            data: data.data.list === undefined || data.data.list === null || data.data.list.length <= 0 ? [] : data.data.list,
            success: data.code === 200,
            total: data.data.total,
          };
        }}
        search={{
          defaultCollapsed: true,
          labelWidth: 'auto',
        }}
        toolBarRender={() => [
          <Button type={"primary"} key="upload_button" onClick={() => {
            setIsModalUploadVisible(true);
          }}>
            上传
          </Button>,
        ]}/>
      {selectedRowsState?.length > 0 && (
        <FooterToolbar
          extra={
            <div>
              已选择&nbsp;&nbsp;<a style={{fontWeight: 600}}>{selectedRowsState.length}</a>&nbsp;&nbsp;项
            </div>
          }>
          <Button danger={true} onClick={modalDelete} hidden={controls["system.file:delete"] == undefined}>
            批量删除
          </Button>
        </FooterToolbar>
      )}
    </PageContainer>
  );
};
