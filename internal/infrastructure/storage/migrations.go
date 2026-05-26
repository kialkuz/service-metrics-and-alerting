package storage

import "github.com/pressly/goose/v3"

func RunMigrations(driver string, dsn string) error {
	db, err := goose.OpenDBWithDriver(driver, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	return goose.Up(db, "migrations")
}
