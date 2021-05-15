package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PSauerborn/gamma-project/internal/pkg/jobs"
	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type PostgresPersistence struct {
	*utils.BasePostgresPersistence
}

// db function used to retrieve a job with given job ID
func (db *PostgresPersistence) GetJob(jobId uuid.UUID) (jobs.Job, error) {
	log.Debug(fmt.Sprintf("fetching job with ID %s", jobId))
	var (
		j    jobs.Job
		meta []byte
	)

	query := `SELECT name,due,meta,state,created,assigned FROM
	jobs WHERE id=$1`
	// get data from database and read into local variables
	row := db.Session.QueryRow(context.Background(), query, jobId)
	if err := row.Scan(&j.Name, &j.Due, &meta, &j.State,
		&j.Created, &j.Assigned); err != nil {
		log.Error(fmt.Errorf("unable to scan data into local variables: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return j, jobs.ErrJobDoesNotExists
		default:
			return j, err
		}
	}
	// convert metadata into JSON and add to struct
	if err := json.Unmarshal(meta, &j.Meta); err != nil {
		log.Error(fmt.Errorf("unable to parse JSON metadata: %+v", err))
		return j, err
	}
	j.JobId = jobId
	return j, nil
}

// db function used to list a collection of jobs
func (db *PostgresPersistence) ListJobs() ([]jobs.Job, error) {
	log.Debug("fetching all jobs from database...")
	results := []jobs.Job{}

	query := `SELECT id,name,due,meta,state,created,assigned FROM jobs`
	rows, err := db.Session.Query(context.Background(), query)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve data from database: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return results, nil
		default:
			return results, err
		}
	}

	for rows.Next() {
		var (
			j    jobs.Job
			meta []byte
		)
		if err := rows.Scan(&j.JobId, &j.Name, &j.Due, &meta, &j.State,
			&j.Created, &j.Assigned); err != nil {
			log.Error(fmt.Errorf("unable to scan data into local variables: %+v", err))
			continue
		}
		// convert metadata into JSON and add to struct
		if err := json.Unmarshal(meta, &j.Meta); err != nil {
			log.Error(fmt.Errorf("unable to parse JSON metadata: %+v", err))
			continue
		}
		results = append(results, j)
	}
	return results, nil
}

func (db *PostgresPersistence) ListUserJobs(uid string) ([]jobs.Job, error) {
	log.Debug(fmt.Sprintf("listing jobs for user %s...", uid))
	results := []jobs.Job{}

	query := `SELECT j.id,j.name,j.due,j.meta,j.state,j.created,j.assigned FROM jobs j
	INNER JOIN assigned_jobs a ON a.id = j.id WHERE a.uid=$1`
	rows, err := db.Session.Query(context.Background(), query, uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to retrieve data from database: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return results, nil
		default:
			return results, err
		}
	}

	for rows.Next() {
		var (
			j    jobs.Job
			meta []byte
		)
		if err := rows.Scan(&j.JobId, &j.Name, &j.Due, &meta, &j.State,
			&j.Created, &j.Assigned); err != nil {
			log.Error(fmt.Errorf("unable to scan data into local variables: %+v", err))
			continue
		}
		// convert metadata into JSON and add to struct
		if err := json.Unmarshal(meta, &j.Meta); err != nil {
			log.Error(fmt.Errorf("unable to parse JSON metadata: %+v", err))
			continue
		}
		results = append(results, j)
	}
	return results, nil
}

func (db *PostgresPersistence) UpdateJobMeta(jobId uuid.UUID, meta map[string]interface{}) error {
	log.Debug(fmt.Sprintf("updating metadata for %s with %+v...", jobId, meta))
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert metadata to JSON: %+v", err))
		return err
	}

	query := `UPDATE jobs SET meta=$1 WHERE id=$2`
	_, err = db.Session.Exec(context.Background(), query, metaJSON, jobId)
	return err
}

// db function used to create a new job
func (db *PostgresPersistence) CreateJob(j jobs.Job) (uuid.UUID, error) {
	log.Debug(fmt.Sprintf("creating new job with values %+v", j))

	id := uuid.New()
	meta, err := json.Marshal(j.Meta)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert metadata to JSON: %+v", err))
		return id, err
	}

	query := `INSERT INTO jobs(id,name,due,meta,state,created)
	VALUES($1,$2,$3,$4,$5,$6)`
	_, err = db.Session.Exec(context.Background(), query, id, j.Name, j.Due, meta, jobs.Created,
		time.Now().UTC())
	if err != nil {
		log.Error(fmt.Errorf("unable to insert job into database: %+v", err))
		return id, err
	}
	return id, nil
}

// db function used to delete a job
func (db *PostgresPersistence) DeleteJob(jobId uuid.UUID) error {
	log.Warn(fmt.Sprintf("deleting job with ID %+v", jobId))
	query := `DELETE FROM jobs WHERE id=$1`
	_, err := db.Session.Exec(context.Background(), query, jobId)
	return err
}

// db function used to alter a job state
func (db *PostgresPersistence) AlterJobState(jobId uuid.UUID, state int) error {
	log.Info(fmt.Sprintf("updating job %s with state %d...", jobId, state))
	query := `UPDATE jobs SET state=$1 WHERE id=$2`
	_, err := db.Session.Exec(context.Background(), query, state, jobId)
	return err
}

// db function to assign jobs to a given user
func (db *PostgresPersistence) AssignJob(jobId uuid.UUID, uid string) error {
	log.Info(fmt.Sprintf("updating job %s with state %s...", jobId, uid))
	var query string
	query = `UPDATE jobs SET state=$1, assigned=true WHERE id=$2`
	_, err := db.Session.Exec(context.Background(), query, jobs.Assigned, jobId)
	if err != nil {
		log.Error(fmt.Errorf("unable to modify job state: %+v", err))
		return err
	}

	query = `INSERT INTO assigned_jobs(id, uid) VALUES($1,$2)
	ON CONFLICT (id) DO
		UPDATE SET uid=$2`
	_, err = db.Session.Exec(context.Background(), query, jobId, uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to assign job: %+v", err))
		return err
	}
	return nil
}
