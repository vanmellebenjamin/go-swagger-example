package services

import (
	"io"
	"log"
	"os"
)

type LocalFileService struct {
	store map[string]string
}

func NewLocalFileService() *LocalFileService {
	fileService := new(LocalFileService)
	fileService.store = make(map[string]string)
	return fileService
}

func (localFileService *LocalFileService) UploadFile(reader io.Reader) (written int64, err error) {
	written, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (localFileService *LocalFileService) DownloadFile() {

}

func (localFileService *LocalFileService) deleteFile() {

}