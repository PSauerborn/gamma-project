package jobs

import (
	"github.com/PSauerborn/gamma-project/internal/pkg/jobs"
	db "github.com/PSauerborn/gamma-project/internal/pkg/jobs/persistence"
	"github.com/PSauerborn/gamma-project/internal/pkg/roles"
	"github.com/PSauerborn/gamma-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

func NewServiceConfig(fileHost, rolesHost string) jobs.ServiceConfig {
	return jobs.ServiceConfig{
		FilestoreHost: fileHost,
		RolesAPIHost:  rolesHost,
	}
}

// function used to generate new instance of jobs API
func NewJobsAPI(p jobs.Persistence, cfg jobs.ServiceConfig) *gin.Engine {
	// set persistence as global within module
	jobs.SetPersistence(p)
	jobs.SetConfig(cfg)
	// generate new instance of gin router and assign routes
	r := gin.Default()
	r.Use(utils.UserHeaderMiddleware())

	r.GET("/jobs/health_check", jobs.HealthCheckHandler)
	// add request handlers to retrieve jobs
	r.GET("/jobs/list/all", utils.RoleMiddelware(roles.Planner, cfg.RolesAPIHost),
		jobs.ListJobsHandler)
	r.GET("/jobs/list", jobs.ListUserJobsHandler)
	r.GET("/jobs/:jobId", jobs.GetJobHandler)

	// add request handler to create new jobs
	r.POST("/jobs/new", utils.RoleMiddelware(roles.Clerk, cfg.RolesAPIHost),
		jobs.CreateJobHandler)
	r.POST("/jobs/:jobId/attachments", jobs.AddJobAttachmentHandler)
	// add request handlers to modify existing jobs
	r.PATCH("/jobs/:jobId/state", jobs.AlterJobStateHandler)
	r.PATCH("/jobs/:jobId/assign", utils.RoleMiddelware(roles.Planner, cfg.RolesAPIHost),
		jobs.AssignJobHandler)
	r.PATCH("/jobs/:jobId/meta", jobs.PatchJobMetaHandler)
	r.DELETE("/jobs/:jobId", utils.RoleMiddelware(roles.Admin, cfg.RolesAPIHost),
		jobs.DeleteJobHandler)
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
