package oss

type ObjectClient interface {
	Upload(object Object, localFilePath string) *UploadRes
	Delete(key string) error
}

type Object struct {
	Name        string
	ContentType string
}

type UploadRes struct {
	Key  string
	Info any
}
