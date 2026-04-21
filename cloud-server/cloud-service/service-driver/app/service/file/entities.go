package file

const (
	ChunkUploadStepInit   = "init"   // 初始化上传步骤
	ChunkUploadStepUpload = "upload" // 上传分片步骤
	ChunkUploadStepStatus = "status" // 查询上传状态步骤
	ChunkUploadStepFinish = "finish" // 完成上传步骤
)

// ChunkUploadResult 分片上传结果。
type ChunkUploadResult struct {
	UploadId       string `json:"upload_id,omitempty"`       // 上传任务ID
	Done           bool   `json:"done,omitempty"`            // 是否完成上传
	Url            string `json:"url,omitempty"`             // 预览地址
	UploadedParts  []int  `json:"uploaded_parts,omitempty"`  // 已上传分片
	CommittedBytes int64  `json:"committed_bytes,omitempty"` // 已提交字节数
}
