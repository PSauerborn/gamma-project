package jobs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrJobDoesNotExists = errors.New("cannot find job with specified ID")

type Persistence interface {
	GetJob(jobId uuid.UUID) (Job, error)
	ListJobs() ([]Job, error)
	CreateJob(job Job) (uuid.UUID, error)
	AssignJob(jobId uuid.UUID, uid string) error
	AlterJobState(jobId uuid.UUID, state int) error
	UpdateJobMeta(jobId uuid.UUID, meta map[string]interface{}) error
	DeleteJob(jobId uuid.UUID) error
}

// generate new type to store job states as enum intergers
type JobState int

// function used to convert job state into a string representation
func (t JobState) String() string {
	return [...]string{"Created", "Assigned", "Completed", "Overdue"}[t]
}

const (
	Created JobState = iota
	Assigned
	Completed
	Overdue
)

type Job struct {
	Name     string                 `json:"name" binding:"required"`
	Due      time.Time              `json:"due" binding:"required"`
	Meta     map[string]interface{} `json:"meta" binding:"required"`
	JobId    uuid.UUID              `json:"job_id"`
	State    JobState               `json:"state"`
	Created  time.Time              `json:"created"`
	Assigned bool                   `json:"assigned"`
}
