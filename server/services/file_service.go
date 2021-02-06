package services

import "io"

type FileService interface {
	UploadFile(reader io.Reader) (written int64, err error)
	DownloadFile()
	deleteFile()
}