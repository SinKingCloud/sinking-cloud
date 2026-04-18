package storage

import (
	"io"
	"time"
)

type uploadProtocol string

const (
	uploadProtocolMultipart uploadProtocol = "multipart"
	uploadProtocolBlock     uploadProtocol = "block"
	uploadProtocolResumable uploadProtocol = "resumable"
)

// chunkPolicy 描述不同厂商对分片大小和数量的约束。
type chunkPolicy struct {
	defaultSize  int64 // 是未显式指定分片大小时采用的默认值。
	minSize      int64 // 是允许的最小分片大小。
	sizeMultiple int64 // 是分片大小必须满足的倍数约束。
	maxChunks    int   // 是单个上传任务允许的最大分片数。
}

// capabilities 描述当前驱动支持的能力边界。
type capabilities struct {
	smallUpload       bool             // 表示是否支持普通直传。
	largeUpload       bool             // 表示是否支持服务端发起的分片上传。
	resumeUpload      bool             // 表示是否支持上传任务继续续传。
	uploadStatusQuery bool             // 表示是否支持查询已上传分片状态。
	signedURL         bool             // 表示是否支持生成带签名的临时访问地址。
	copyObject        bool             // 表示是否支持服务端对象复制。
	batchDelete       bool             // 表示是否支持批量删除对象。
	protocols         []uploadProtocol // 描述驱动支持的上传协议模型。
	chunkPolicy       chunkPolicy      // 描述当前驱动使用的分片策略。
}

// uploadRequest 描述一次普通上传所需的内部参数。
type uploadRequest struct {
	key                string            // 是业务层传入的逻辑对象路径。
	body               io.Reader         // 是待上传的内容流。
	size               int64             // 是对象总大小，单位为字节。
	contentType        string            // 是对象的内容类型。
	cacheControl       string            // 是对象的缓存控制头。
	contentDisposition string            // 是对象的下载展示方式。
	metadata           map[string]string // 是对象自定义元数据。
}

// largeUploadRequest 描述一次自动分片上传所需的内部参数。
type largeUploadRequest struct {
	uploadRequest       // 复用普通上传的基础参数。
	chunkSize     int64 // 是期望使用的分片大小。
}

// startUploadSessionRequest 描述初始化上传会话所需的内部参数。
type startUploadSessionRequest struct {
	key                string            // 是业务层传入的逻辑对象路径。
	size               int64             // 是整个文件总大小，单位为字节。
	chunkSize          int64             // 是期望采用的分片大小。
	contentType        string            // 是对象的内容类型。
	cacheControl       string            // 是对象的缓存控制头。
	contentDisposition string            // 是对象的下载展示方式。
	metadata           map[string]string // 是对象自定义元数据。
}

// uploadSession 描述一个可继续续传的内部上传会话。
type uploadSession struct {
	provider           Provider          // 标识当前会话所属的存储厂商。
	protocol           uploadProtocol    // 标识当前会话采用的上传协议。
	bucket             string            // 是当前会话所属的 Bucket 名称。
	key                string            // 是业务层看到的逻辑对象路径。
	providerKey        string            // 是底层厂商实际使用的对象路径。
	sessionID          string            // 是厂商返回的原始上传会话 ID。
	size               int64             // 是整个文件总大小，单位为字节。
	chunkSize          int64             // 是当前会话采用的分片大小。
	contentType        string            // 是对象的内容类型。
	cacheControl       string            // 是对象的缓存控制头。
	contentDisposition string            // 是对象的下载展示方式。
	metadata           map[string]string // 是对象自定义元数据。
	opaque             map[string]string // 用于承载不同厂商私有扩展字段。
	createdAt          time.Time         // 是会话初始化时间。
}

// uploadChunkRequest 描述单个分片上传所需的内部参数。
type uploadChunkRequest struct {
	session  uploadSession // 是当前分片所属的上传会话。
	index    int           // 是分片序号，通常从 1 开始。
	offset   int64         // 是分片在整个文件中的起始偏移。
	size     int64         // 是当前分片的大小，单位为字节。
	chunkID  string        // 是用于兼容不同厂商的分片标识。
	isFinal  bool          // 表示当前分片是否为最后一片。
	body     io.Reader     // 是当前分片的内容流。
	checksum string        // 是当前分片的可选校验值。
}

// uploadedChunk 描述云端已经成功接收的分片信息。
type uploadedChunk struct {
	index        int               // 是分片序号。
	offset       int64             // 是分片在整个文件中的起始偏移。
	size         int64             // 是分片大小，单位为字节。
	chunkID      string            // 是兼容不同厂商时使用的分片标识。
	etag         string            // 是厂商返回的分片内容标识。
	checksum     string            // 是分片的校验结果。
	lastModified time.Time         // 是分片最近一次成功上传时间。
	opaque       map[string]string // 用于承载不同厂商的私有扩展字段。
}

// uploadStatus 描述一个上传会话当前的内部状态。
type uploadStatus struct {
	session        uploadSession   // 是当前查询对应的上传会话。
	uploadedChunks []uploadedChunk // 是云端已经确认接收的分片列表。
	committedBytes int64           // 是云端已确认写入的总字节数。
	finalized      bool            // 表示上传任务是否已经完成。
}

// completeUploadRequest 描述完成上传时提交给驱动的内部参数。
type completeUploadRequest struct {
	session uploadSession   // 是待完成的上传会话。
	chunks  []uploadedChunk // 是完成上传所需的分片清单。
}
