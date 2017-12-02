package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"

	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/migrations"
)

var databaseURL string

func init() {
	// Register command under the root command
	rootCmd.AddCommand(migrateCmd)

	// Register subcommands
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateShowCmd)

	// Define persistent flags for all commands under the migrateCmd
	migrateCmd.PersistentFlags().StringVarP(
		&databaseURL,
		"database-url", "d", "",
		"database source name")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage schema and data migrations",
}

// TODO: Decide if we need a shortcut to create up&down migration files

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrater, migraterErr := newMigrate()
		defer func() {
			if migraterErr == nil {
				migrater.Close()
			}
		}()

		if err := migrater.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatalln("error:", err)
			} else {
				log.Println("error:", err)
			}
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Apply down migration",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Not implemented")
	},
}

var migrateShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current status of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Not implemented")
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
		"go-bindata", bindataMigrateSource, databaseURL)

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
