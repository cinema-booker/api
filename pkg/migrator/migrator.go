package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"time"
)

type Migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(db *sql.DB, dir string) (*Migrator, error) {
	// get the relative path of the migrations directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// create a migrator instance
	m := &Migrator{
		db:  db,
		dir: path.Join(cwd, "migrations"),
	}
	return m, nil
}

func (m *Migrator) CreateMigration(migrationName string) error {
	// check if the migrations directory exists
	if _, err := os.Stat(m.dir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("migrations directory does not exist: %v", err)
		}
		return fmt.Errorf("error getting migrations directory: %v", err)
	}

	// generate new migration version: YYYY_MM_DD_HHMMSS
	version := time.Now().Format("2006_01_02_150405")

	// create the migration `up` file
	upFileName := fmt.Sprintf("%s_%s.up.sql", version, migrationName)
	upFilePath := path.Join(m.dir, upFileName)
	if _, err := os.Stat(upFilePath); err == nil {
		return fmt.Errorf("migration file with name '%s' already exists", upFileName)
	}
	upFile, err := os.Create(upFilePath)
	if err != nil {
		return fmt.Errorf("error creating up migration file: %v", err)
	}
	defer upFile.Close()

	// create the migration `down` file
	downFileName := fmt.Sprintf("%s_%s.down.sql", version, migrationName)
	downFilePath := path.Join(m.dir, downFileName)
	if _, err := os.Stat(downFilePath); err == nil {
		return fmt.Errorf("migration file with name '%s' already exists", downFileName)
	}
	downFile, err := os.Create(downFilePath)
	if err != nil {
		return fmt.Errorf("error creating down migration file: %v", err)
	}
	defer downFile.Close()

	return nil
}
