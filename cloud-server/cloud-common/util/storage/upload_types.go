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

type chunkPolicy struct {
	defaultSize  int64
	minSize      int64
	sizeMultiple int64
	maxChunks    int
}

type capabilities struct {
	smallUpload       bool
	largeUpload       bool
	resumeUpload      bool
	uploadStatusQuery bool
	signedURL         bool
	copyObject        bool
	batchDelete       bool
	protocols         []uploadProtocol
	chunkPolicy       chunkPolicy
}

type uploadRequest struct {
	key                string
	body               io.Reader
	size               int64
	contentType        string
	cacheControl       string
	contentDisposition string
	metadata           map[string]string
}

type largeUploadRequest struct {
	uploadRequest
	chunkSize int64
}

type startUploadSessionRequest struct {
	key                string
	size               int64
	chunkSize          int64
	contentType        string
	cacheControl       string
	contentDisposition string
	metadata           map[string]string
}

type uploadSession struct {
	provider           Provider
	protocol           uploadProtocol
	bucket             string
	key                string
	providerKey        string
	sessionID          string
	size               int64
	chunkSize          int64
	contentType        string
	cacheControl       string
	contentDisposition string
	metadata           map[string]string
	opaque             map[string]string
	createdAt          time.Time
}

type uploadChunkRequest struct {
	session  uploadSession
	index    int
	offset   int64
	size     int64
	chunkID  string
	isFinal  bool
	body     io.Reader
	checksum string
}

type uploadedChunk struct {
	index        int
	offset       int64
	size         int64
	chunkID      string
	etag         string
	checksum     string
	lastModified time.Time
	opaque       map[string]string
}

type uploadStatus struct {
	session        uploadSession
	uploadedChunks []uploadedChunk
	committedBytes int64
	finalized      bool
}

type completeUploadRequest struct {
	session uploadSession
	chunks  []uploadedChunk
}
