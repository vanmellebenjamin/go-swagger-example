package services

import (
	"flightAPI/server/models"
	"github.com/go-openapi/strfmt"
	"io"
	"os"
)

type FileService interface {
	UploadFile(reader io.Reader) (*models.File, error)
	DownloadFile(uuid strfmt.UUID) (*os.File, error)
	DeleteFile(uuid strfmt.UUID) error
}
