package justcheck

import (
	"context"
	"mime/multipart"
)

type UseCase interface {
	Check(ctx context.Context,file multipart.File,contentType string) (string,error)
}