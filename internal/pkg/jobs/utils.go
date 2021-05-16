package jobs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFileStorageError = errors.New("unable to store file in filestore")
	ErrInvalidJobID     = errors.New("received invalid job ID")
)

func ParseAndValidateJobId(ctx *gin.Context, key string) (uuid.UUID, error) {
	id, err := uuid.Parse(ctx.Param(key))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job id: %+v", err))
		return id, ErrInvalidJobID
	}

	_, err = persistence.GetJob(id)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		return id, err
	}
	return id, nil
}

// function used to update job metadata in database via JSON patch operation
func UpdateJobMetadata(jobId uuid.UUID, patch []map[string]interface{}) error {
	log.Debug(fmt.Sprintf("patching metadata for job %+v", jobId))
	job, err := persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		return err
	}
	// perform JSON patch operation on metadata
	patched, err := utils.PatchJSON(job.Meta, patch)
	if err != nil {
		log.Error(fmt.Errorf("unable to perform JSON patch: %+v", err))
		return err
	}
	return persistence.UpdateJobMeta(jobId, patched)
}

// function used to append an attachment ID to a list of
// attachments
func AddJobAttachment(jobId, fileId uuid.UUID) error {
	log.Debug(fmt.Sprintf("adding file %s to job %s", fileId, jobId))
	job, err := persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		return err
	}
	// get attachments and convert to string slice
	attachments := job.Meta["attachments"].([]interface{})
	job.Meta["attachments"] = append(attachments, fileId.String())
	return persistence.UpdateJobMeta(jobId, job.Meta)
}

func AddFileToFilestore(upload FileUpload, host string) (uuid.UUID, error) {
	log.Info("adding new file to filestore")

	var r struct {
		HTTPCode int       `json:"http_code"`
		FileId   uuid.UUID `json:"file_id"`
	}

	jsonBody, err := json.Marshal(upload)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert file to JSON format: %+v", err))
		return r.FileId, err
	}

	url := fmt.Sprintf("%s/filestore/file", host)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error(fmt.Errorf("unable to generate HTTP request: %+v", err))
		return r.FileId, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Authenticated-Userid", "jobs-api")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute HTTP request: %+v", err))
		return r.FileId, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		log.Info("successfully added file to filestore")
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			log.Error(fmt.Errorf("unable to parse API response: %+v", err))
			return r.FileId, err
		}
		return r.FileId, nil
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("failed to read API response: %+v", err)
		} else {
			log.Error(fmt.Sprintf("received invalid API response with code %d: %s",
				resp.StatusCode, string(body)))
		}
		return r.FileId, ErrFileStorageError
	}
}
