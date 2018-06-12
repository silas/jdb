package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const version = 1

type Revision struct {
	id  int
	sql string
}

func SQL(id int, sql string) Revision {
	return Revision{id: id, sql: sql}
}

type Revisions []Revision

func (revisions Revisions) Run(ctx context.Context, db *sql.DB, helper Helper) error {
	if ctx == nil {
		return errors.New("ctx required")
	}
	if db == nil {
		return errors.New("db required")
	}
	if helper == nil {
		return errors.New("helper required")
	}

	if helper.Version() != version {
		return fmt.Errorf("dialect version %d does not match jdb version %d", helper.Version(), version)
	}

	currentID, err := getRevisionID(ctx, db, helper)
	if err != nil {
		return err
	}

	if currentID > 0 {
		dbVersion, err := getVersion(ctx, db, helper)
		if err != nil {
			return err
		}
		if dbVersion > version {
			return fmt.Errorf("database version %d is newer than library version %d", dbVersion, version)
		}
	}

	for i, r := range revisions {
		id := i + 1

		if r.id != r.id {
			return fmt.Errorf("invalid revision ID: %d != %d", r.id, id)
		}

		if id <= currentID {
			continue
		}

		if err := run(ctx, db, helper, r); err != nil {
			return err
		}
	}

	return setVersion(ctx, db, helper)
}

func getRevisionID(ctx context.Context, db *sql.DB, helper Helper) (int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if exists, err := helper.TableExists(ctx, tx); err != nil {
		return 0, err
	} else if !exists {
		return 0, nil
	}

	return helper.GetRevision(ctx, tx)
}

func getVersion(ctx context.Context, db *sql.DB, helper Helper) (int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	return helper.GetVersion(ctx, tx)
}

func setVersion(ctx context.Context, db *sql.DB, helper Helper) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = helper.SetVersion(ctx, tx, version)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func run(ctx context.Context, db *sql.DB, helper Helper, revision Revision) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	local := map[string]interface{}{
		"ID": revision.id,
	}
	_, err = tx.ExecContext(ctx, helper.RenderSQL(revision.sql, local))
	if err != nil {
		return err
	}

	err = helper.SetRevision(ctx, tx, revision.id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
