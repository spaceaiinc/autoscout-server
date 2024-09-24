package tests

import (
	"github.com/spaceaiinc/autoscout-server/domain/config"
	"github.com/spaceaiinc/autoscout-server/infrastructure/database"
)

type Helper struct {
	dbm *database.DB
}

func NewHelper(cfg config.Config) *Helper {
	dbm := database.NewDB(cfg.DB, false)
	err := dbm.MigrateUp("../../.migrations")
	if err != nil {
		panic(err)
	}

	return &Helper{
		dbm: dbm,
	}
}

func (h *Helper) ClearAllTables() {
	h.dbm.Exec(
		`DELETE FROM users`,
		`DROP TABLE users`,
	)
}
