package filestore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
)

type FileStoreAPIAccessor struct {
	*utils.BaseAPIAccessor
}

func (accessor *FileStoreAPIAccessor) GetFileMetadata(fileId uuid.UUID) (FileMetadata, error) {
	log.Debug(fmt.Sprintf("retrieving file metadata for %s", fileId))
	var payload struct {
		HTTPCode int          `json:"http_code"`
		Metadata FileMetadata `json:"metadata"`
	}
	// generate URL using file ID, and create new request
	url := accessor.FormatURL(fmt.Sprintf("/filestore/file/meta/%s", fileId))
	request, err := accessor.NewJSONRequest("GET", url, nil, nil)
	if err != nil {
		log.Error(fmt.Errorf("unable to generate new request: %+v", err))
		return payload.Metadata, err
	}

	response, err := accessor.ExecuteRequest(request)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute request: %+v", err))
		return payload.Metadata, err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 200:
		log.Debug("successfully retrieved metadata from API")
		if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
			log.Error(fmt.Errorf("unable to parse JSON response from API: %+v", err))
			return payload.Metadata, err
		}
		return payload.Metadata, nil
	default:
		body, _ := ioutil.ReadAll(response.Body)
		log.Error(fmt.Sprintf("received non-success response from API: %+v", string(body)))
	}
	return payload.Metadata, nil
}

func (accessor *FileStoreAPIAccessor) GetFileContents(fileId uuid.UUID) ([]byte, error) {
	var contents []byte
	return contents, nil
}

func (accessor *FileStoreAPIAccessor) ListFiles() ([]FileMetadata, error) {
	files := []FileMetadata{}
	return files, nil
}

func (accessor *FileStoreAPIAccessor) CreateFile(contents []byte) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (accessor *FileStoreAPIAccessor) ModifyFile(fileId uuid.UUID) error {
	return nil
}

func (accessor *FileStoreAPIAccessor) DeleteFile(fileId uuid.UUID) error {
	return nil
}
