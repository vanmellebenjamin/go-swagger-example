package services

import (
	"errors"
	"flightAPI/server/models"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
)

const fileDirectory = "C:\\Users\\vanme\\go\\temp\\%s"

type LocalFileService struct {
	store map[string]string
}

func NewLocalFileService() *LocalFileService {
	fileService := new(LocalFileService)
	fileService.store = make(map[string]string)
	return fileService
}

func (localFileService *LocalFileService) UploadFile(reader io.Reader) (*models.File, error) {
	id := strfmt.UUID(uuid.NewString())
	file := models.File{UUID: &id}
	tgtFile, err := os.Create(fmt.Sprintf(fileDirectory, id))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer tgtFile.Close()
	if _, err = io.Copy(tgtFile, reader); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &file, nil
}

func (localFileService *LocalFileService) DownloadFile(uuid strfmt.UUID) (*os.File, error) {
	if !fileExists(uuid) {
		return nil, errors.New("not_found")
	}
	file, err := os.Open(fmt.Sprintf(fileDirectory, uuid))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return file, nil
}

func (localFileService *LocalFileService) DeleteFile(uuid strfmt.UUID) error {
	if !fileExists(uuid) {
		return errors.New("not_found")
	}
	if err := os.Remove(fmt.Sprintf(fileDirectory, uuid)); err != nil {
		return err
	}
	return nil
}

func fileExists(uuid strfmt.UUID) bool {
	info, err := os.Stat(fmt.Sprintf(fileDirectory, uuid))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}