package jobs

import (
	"github.com/PSauerborn/gamma-project/internal/pkg/jobs"
	db "github.com/PSauerborn/gamma-project/internal/pkg/jobs/persistence"
	"github.com/PSauerborn/gamma-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

// function used to generate new instance of jobs API
func NewJobsAPI(p jobs.Persistence) *gin.Engine {
	// set persistence as global within module
	jobs.SetPersistence(p)
	// generate new instance of gin router and assign routes
	r := gin.Default()
	r.Use(utils.UserHeaderMiddleware())

	r.GET("/jobs/health_check", jobs.HealthCheckHandler)
	// add request handlers to retrieve jobs
	r.GET("/jobs/list/all", jobs.ListJobsHandler)
	r.GET("/jobs/list", jobs.ListUserJobsHandler)
	r.GET("/jobs/:jobId", jobs.GetJobHandler)

	// add request handler to create new jobs
	r.POST("/jobs/new", jobs.CreateJobHandler)
	// add request handlers to modify existing jobs
	r.PATCH("/jobs/:jobId/state", jobs.AlterJobStateHandler)
	r.PATCH("/jobs/:jobId/assign", jobs.AssignJobHandler)
	r.PATCH("/jobs/:jobId/meta", jobs.PatchJobMetaHandler)
	r.DELETE("/jobs/:jobId", jobs.DeleteJobHandler)
	return r
}

// function used to generate new instance of postgres persistence
func NewPostgresPersistence(url string) *db.PostgresPersistence {
	// generate base persistence layer
	base := utils.NewBasePersistence(url)
	return &db.PostgresPersistence{
		BasePostgresPersistence: base,
	}
}
