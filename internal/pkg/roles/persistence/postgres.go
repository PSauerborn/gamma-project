package roles

import (
	"context"
	"fmt"

	"github.com/PSauerborn/gamma-project/internal/pkg/roles"
	"github.com/PSauerborn/gamma-project/internal/pkg/utils"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type PostgresPersistence struct {
	*utils.BasePostgresPersistence
}

func (db *PostgresPersistence) GetUserRole(uid string) (roles.Role, error) {
	log.Debug(fmt.Sprintf("feching role for user %s...", uid))
	var r roles.Role
	query := `SELECT role FROM user_roles WHERE uid=$1`
	row := db.Session.QueryRow(context.Background(), query, uid)
	if err := row.Scan(&r); err != nil {
		log.Error(fmt.Errorf("unable to scan data into local variables: %+v", err))
		switch err {
		case pgx.ErrNoRows:
			return roles.StandardUser, nil
		default:
			return r, err
		}
	}
	return r, nil
}

func (db *PostgresPersistence) SetUserRole(uid string, r roles.Role) error {
	log.Debug(fmt.Sprintf("setting user %s with role %d...", uid, r))
	query := `INSERT INTO user_roles(uid, role) VALUES($1,$2)
	ON CONFLICT (uid) DO UPDATE SET role = $2`
	_, err := db.Session.Exec(context.Background(), query, uid, r)
	return err
}
