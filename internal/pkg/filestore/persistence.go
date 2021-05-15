package filestore

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// define custom errors for file persistence
	ErrFileNotFound        = errors.New("cannot find specified file")
	ErrPermissionDenied    = errors.New("permission denied when trying to access file")
	ErrCannotDeleteFile    = errors.New("cannot delete specified file")
	ErrFeatureNotSupported = errors.New("selectd feature currently not supported")
)

var persistence FileStorePersistence

// function used to set global instance of file
// storage persistence
func SetFilePersistence(p FileStorePersistence) {
	persistence = p
}

// define interface for persistence file data. note
// that the files themselves are stored on disk: it
// is merely the file information that is stored in
// a separate persistence layer
type FileStorePersistence interface {
	ListFiles() ([]FileMetadata, error)
	GetFileMetadata(fileId uuid.UUID) (FileMetadata, error)
	GetFileContents(meta FileMetadata) ([]byte, error)
	CreateFile(contents []byte, fileName string, meta map[string]interface{}) (uuid.UUID, error)
	ModifyFile(meta FileMetadata, contents []byte) error
	DeleteFile(meta FileMetadata) error
	ArchiveFile(meta FileMetadata) error
	SearchFilesByMetadata(terms map[string]interface{}) ([]FileMetadata, error)
}

type FileMetadata struct {
	FileId   uuid.UUID              `json:"file_id" validate:"required"`
	FileName string                 `json:"file_name" validate:"required"`
	Created  time.Time              `json:"created" validate:"required"`
	Size     int                    `json:"size" validate:"required"`
	Meta     map[string]interface{} `json:"meta" validate:"required"`
}
