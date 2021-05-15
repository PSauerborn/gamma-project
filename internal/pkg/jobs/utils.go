package jobs

import (
	"fmt"

	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// function used to update job metadata in database via JSON patch operation
func UpdateJobMetadata(jobId uuid.UUID, patch []map[string]interface{}) error {
	log.Debug(fmt.Sprintf("patching metadata for job %+v", jobId))
	job, err := persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		return err
	}

	patched, err := utils.PatchJSON(job.Meta, patch)
	if err != nil {
		log.Error(fmt.Errorf("unable to perform JSON patch: %+v", err))
		return err
	}
	return persistence.UpdateJobMeta(jobId, patched)
}
