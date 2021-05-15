package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type BasePostgresPersistence struct {
	DatabaseURL string
	Session     *pgxpool.Pool
}

// function to connect persistence to postgres server
// note that the connection is returned and should be
// closed with a defer conn.Close(context) statement
func (db *BasePostgresPersistence) Connect() error {
	log.Debug("creating new database connection")
	// connect to postgres server and set session in persistence
	conn, err := pgxpool.Connect(context.Background(), db.DatabaseURL)
	if err != nil {
		log.Error(fmt.Errorf("error connecting to postgres service: %+v", err))
		return err
	}
	db.Session = conn
	return err
}

func (db *BasePostgresPersistence) Close() {
	db.Session.Close()
}
