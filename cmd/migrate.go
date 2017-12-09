package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mattes/migrate"
	// Init migrate's postgres driver
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"

	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/env"
	"github.com/m1kola/shipsterbot/migrations"
)

func init() {
	// Register command under the root command
	rootCmd.AddCommand(migrateCmd)

	// Register subcommands
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateShowCurrentCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage schema and data migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m, migraterErr := newMigrate()
		defer migrateCleanup(m, migraterErr)

		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Nothing to migrate")
			} else {
				log.Fatalln("error:", err)
			}
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply down migration",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg, received %d", len(args))
		}

		if args[0] == "zero" {
			return nil
		}

		version, err := strconv.Atoi(args[0])
		if err != nil || version < 1 {
			return fmt.Errorf(
				"accepts a positive int or \"zero\" to apply all down migrations, received %s",
				args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		m, migraterErr := newMigrate()
		defer migrateCleanup(m, migraterErr)

		if args[0] == "zero" {
			// If arg is "zero" - apply all down migrations

			err = m.Down()
		} else {
			// Othervise migrate to a specific version

			var version int
			// No need to check error here: we've already
			// validated the argument in our
			// custom Args validator function
			version, _ = strconv.Atoi(args[0])
			err = m.Migrate(uint(version))
		}

		if err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Nothing to migrate")
			} else {
				log.Fatalln("error:", err)
			}
		}
	},
}

var migrateShowCurrentCmd = &cobra.Command{
	Use:   "showcurrent",
	Short: "Shows current migration's number",
	Run: func(cmd *cobra.Command, args []string) {
		m, migraterErr := newMigrate()
		defer migrateCleanup(m, migraterErr)

		currentVersion, dirty, err := m.Version()
		if err != nil {
			log.Fatalln("error:", err)
		}

		if dirty {
			fmt.Printf("Current migration: \"%v\" (dirty)\n", currentVersion)
		} else {
			fmt.Printf("Current migration: %v\n", currentVersion)
		}
	},
}

func newMigrate() (*migrate.Migrate, error) {
	// Initialize bindata resources
	// Doesn't make much sense to be configurable, at the momnet
	bindataResource := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	bindataMigrateSource, err := bindata.WithInstance(bindataResource)
	if err != nil {
		log.Fatalln("error:", err)
	}

	// Initialize migrate
	// Each command must decide  how it wants to handle the migraterErr
	migrater, migraterErr := migrate.NewWithSourceInstance(
		"go-bindata", bindataMigrateSource, env.GetDBConnectionString())

	if migraterErr == nil {
		// handle Ctrl+c
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT)
		go func() {
			for range signals {
				log.Println("Stopping after this running migration ...")
				migrater.GracefulStop <- true
				return
			}
		}()
	}

	if migraterErr != nil {
		log.Fatalln("error:", migraterErr)
	}

	return migrater, migraterErr
}

func migrateCleanup(m *migrate.Migrate, migraterErr error) {
	if migraterErr == nil {
		m.Close()
	}
}
