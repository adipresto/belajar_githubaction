package file

import (
	"context"
	"mime/multipart"
	"os"
)

type UseCase interface {
	UploadFile(context.Context, *multipart.FileHeader, multipart.File, string) (string, error)
	GetFile(context.Context, string) (*os.File, error)
}
