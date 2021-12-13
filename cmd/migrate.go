package cmd

import (
	"ibtools_server/models"
	"ibtools_server/util/migrations"
)

// Migrate runs database migrations
func Migrate(configBackend string) error {
	_, db, _, _, _, err := initConfigDB(true, false, configBackend)
	if err != nil {
		return err
	}

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := models.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
