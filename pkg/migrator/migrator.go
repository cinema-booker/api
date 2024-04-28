package migrator

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
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

func (m *Migrator) MigrateUp(step int) error {
	// start migration transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	// create the migrations table if it does not exist
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS migrations (version varchar(17))"); err != nil {
		return err
	}

	// get the current version from the database
	var version sql.NullString
	if err := tx.QueryRow("SELECT version FROM migrations LIMIT 1").Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			if _, err := tx.Exec("INSERT INTO migrations (version) VALUES (NULL)"); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// get migration files (sorted alphabetically by default)
	files, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}

	// check if the migration files are corrupted
	if version.Valid {
		isCorrupted := true
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), version.String) {
				isCorrupted = false
				break
			}
		}
		if isCorrupted {
			return fmt.Errorf("migration files are corrupted")
		}
	}

	// get the migration files to apply
	var migrations []string
	for i := 0; i <= len(files)-1; i++ {
		file := files[i]
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".up.sql") {
			if !version.Valid {
				migrations = append(migrations, file.Name())
			} else if version.Valid && !strings.HasPrefix(file.Name(), version.String) && file.Name() > version.String {
				migrations = append(migrations, file.Name())
			}
			if len(migrations) == step {
				break
			}
		}
	}
	if len(migrations) == 0 {
		return fmt.Errorf("no migration files to apply")
	}

	// execute the migration files sql queries
	for _, migration := range migrations {
		// get the migration file
		migrationPath := path.Join(m.dir, migration)
		migrationFile, err := os.Open(migrationPath)
		if err != nil {
			return err
		}
		defer migrationFile.Close()
		// read the migration file content
		data, err := io.ReadAll(migrationFile)
		if err != nil {
			return err
		}
		// execute the migration file content
		if _, err = tx.Exec(string(data)); err != nil {
			return err
		}
		// TODO: use logger
		fmt.Printf("Applied migration: %s\n", migration)
	}

	// get the new version after applying the migrations
	fileName := migrations[len(migrations)-1]
	parts := strings.Split(fileName, "_")
	newVersion := strings.Join(parts[:4], "_")

	// save the new version in the database
	if _, err := tx.Exec("UPDATE migrations SET version = $1", newVersion); err != nil {
		return err
	}

	// commit the migration transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) MigrateDown(step int) error {
	// start migration transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	// create the migrations table if it does not exist
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS migrations (version varchar(17))"); err != nil {
		return err
	}

	// get the current version from the database
	var version sql.NullString
	if err := tx.QueryRow("SELECT version FROM migrations LIMIT 1").Scan(&version); err != nil {
		if err == sql.ErrNoRows {
			if _, err := tx.Exec("INSERT INTO migrations (version) VALUES (NULL)"); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// version is NULL, no migration files to rollback
	if !version.Valid {
		return fmt.Errorf("no migration files to rollback")
	}

	// get migration files (sorted alphabetically by default)
	files, err := os.ReadDir(m.dir)
	if err != nil {
		return err
	}

	// check if the migration files are corrupted
	isCorrupted := true
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), version.String) {
			isCorrupted = false
			break
		}
	}
	if isCorrupted {
		return fmt.Errorf("migration files are corrupted")
	}

	// get the migration files to apply
	var migrations []string
	newVersion := sql.NullString{
		String: "",
		Valid:  false,
	}
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".down.sql") {
			if strings.HasPrefix(file.Name(), version.String) || file.Name() < version.String {
				migrations = append(migrations, file.Name())
				if i == 0 || len(migrations) == step {
					if i <= 1 {
						newVersion = sql.NullString{
							String: "",
							Valid:  false,
						}
					} else {
						parts := strings.Split(files[i-1].Name(), "_")
						newVersion = sql.NullString{
							String: strings.Join(parts[:4], "_"),
							Valid:  true,
						}
					}
					break
				}
			}
		}
	}
	if len(migrations) == 0 {
		return fmt.Errorf("no migration files to apply")
	}

	// execute the migration files sql queries
	for _, migration := range migrations {
		// get the migration file
		migrationPath := path.Join(m.dir, migration)
		migrationFile, err := os.Open(migrationPath)
		if err != nil {
			return err
		}
		defer migrationFile.Close()
		// read the migration file content
		data, err := io.ReadAll(migrationFile)
		if err != nil {
			return err
		}
		// execute the migration file content
		if _, err = tx.Exec(string(data)); err != nil {
			return err
		}
		// TODO: use logger
		fmt.Printf("Reverted migration: %s\n", migration)
	}

	// save the new version in the database
	if _, err := tx.Exec("UPDATE migrations SET version = $1", newVersion); err != nil {
		return err
	}

	// commit the migration transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
