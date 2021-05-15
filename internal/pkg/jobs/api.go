package jobs

import (
	"fmt"
	"net/http"

	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var persistence Persistence

// function used to set global persistence instance
func SetPersistence(p Persistence) {
	persistence = p
}

// API handler used to serve health check routes
func HealthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check route")
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Service running"})
}

// API handler used to list all jobs
func ListJobsHandler(ctx *gin.Context) {
	log.Info("received request to list jobs")
	// get all jobs from persistence layer
	jobs, err := persistence.ListJobs()
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve jobs: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"jobs": jobs})
}

// API handler used to list all jobs
func ListUserJobsHandler(ctx *gin.Context) {
	log.Info("received request to list jobs for user")
	uid := ctx.MustGet("uid").(string)
	// get all jobs from persistence layer
	jobs, err := persistence.ListUserJobs(uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve jobs: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"jobs": jobs})
}

// API handler used to retrieve a job with given job ID
func GetJobHandler(ctx *gin.Context) {
	log.Info("received request to retrieve job")
	// extract job ID from path and parse
	jobId, err := uuid.Parse(ctx.Param("jobId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job ID: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid job ID"})
		return
	}
	// get job from persistence layer
	j, err := persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job: %+v", err))
		switch err {
		case ErrJobDoesNotExists:
			status := http.StatusNotFound
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Cannot find job with specified ID"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"job": j})
}

// API handler used to create new jobs
func CreateJobHandler(ctx *gin.Context) {
	log.Info("received request to create new job")
	var j Job
	if err := ctx.ShouldBind(&j); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid request body"})
		return
	}
	// create new job in persistence layer
	id, err := persistence.CreateJob(j)
	if err != nil {
		log.Error(fmt.Errorf("unable to create new job: %+v", err))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
		"message": "Successfully created job", "id": id})
}

// API handler used to delete job
func DeleteJobHandler(ctx *gin.Context) {
	log.Info("received request to delete job")
	// extract job ID from path and parse
	jobId, err := uuid.Parse(ctx.Param("jobId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job ID: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid job ID"})
		return
	}

	// get job details from database
	_, err = persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		switch err {
		case ErrJobDoesNotExists:
			status := http.StatusNotFound
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Cannot find job with specified ID"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	if err := persistence.DeleteJob(jobId); err != nil {
		log.Error(fmt.Errorf("unable to delete job from database"))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully delete job"})
}

// API handler used to alter job state
func AlterJobStateHandler(ctx *gin.Context) {
	log.Info("received request to update job state")
	var r struct {
		State int `json:"state" binding:"required"`
	}
	if err := ctx.ShouldBind(&r); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid request body"})
		return
	}

	// extract job ID from path and parse
	jobId, err := uuid.Parse(ctx.Param("jobId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job ID: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid job ID"})
		return
	}
	// get job details from database
	_, err = persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		switch err {
		case ErrJobDoesNotExists:
			status := http.StatusNotFound
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Cannot find job with specified ID"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	if err := persistence.AlterJobState(jobId, r.State); err != nil {
		log.Error(fmt.Errorf("unable to alter job state"))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully updated job"})
}

// API handler used to assign job
func AssignJobHandler(ctx *gin.Context) {
	log.Info("received request to assign job")
	var r struct {
		User string `json:"user" binding:"required"`
	}
	if err := ctx.ShouldBind(&r); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid request body"})
		return
	}

	// extract job ID from path and parse
	jobId, err := uuid.Parse(ctx.Param("jobId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job ID: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid job ID"})
		return
	}
	// get job details from database
	_, err = persistence.GetJob(jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve job from database: %+v", err))
		switch err {
		case ErrJobDoesNotExists:
			status := http.StatusNotFound
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Cannot find job with specified ID"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	if err := persistence.AssignJob(jobId, r.User); err != nil {
		log.Error(fmt.Errorf("unable to assign job"))
		status := http.StatusInternalServerError
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully updated job"})
}

func PatchJobMetaHandler(ctx *gin.Context) {
	log.Info("received request to patch job metadata")
	// extract job ID from path and parse
	jobId, err := uuid.Parse(ctx.Param("jobId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse job ID: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid job ID"})
		return
	}

	var r struct {
		Operation []map[string]interface{} `json:"operation" binding:"required"`
	}
	if err := ctx.ShouldBind(&r); err != nil {
		log.Error(fmt.Errorf("unable to parse request body: %+v", err))
		status := http.StatusBadRequest
		ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
			"message": "Invalid request body"})
		return
	}

	if err := UpdateJobMetadata(jobId, r.Operation); err != nil {
		log.Error(fmt.Errorf("unable to perform JSON patch: %+v", err))
		switch err {
		case ErrJobDoesNotExists:
			status := http.StatusNotFound
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Cannot find job with specified ID"})
		case utils.ErrInvalidPatch:
			status := http.StatusBadRequest
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Invalid JSON Patch operation"})
		default:
			status := http.StatusInternalServerError
			ctx.AbortWithStatusJSON(status, gin.H{"http_code": status,
				"message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"message": "Successfully patched job metadata"})
}
