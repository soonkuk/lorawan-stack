// Copyright © 2020 The Things Industries B.V.

package commands

import (
	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/v3/pkg/applicationserver/io/packages/storage/postgres"
)

var (
	storageDBCommand = &cobra.Command{
		Use:   "storage-db",
		Short: "Manage the Storage Integration database",
	}
	storageDBInitCommand = &cobra.Command{
		Use:   "init",
		Short: "Initialize the Storage Integration database",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Connecting to Storage Integration database...")

			switch cfg := config.AS.Packages.Config.Storage; cfg.Provider {
			case "postgres":
				db, err := postgres.Open(cfg.Postgres.DatabaseURI, cfg.Postgres)
				if err != nil {
					return err
				}
				defer db.Close()
				if err := postgres.Initialize(db); err != nil {
					return err
				}
			default:
				logger.Info("No database configured")
				return nil
			}

			logger.Info("Successfully initialized")
			return nil
		},
	}
)

func init() {
	Root.AddCommand(storageDBCommand)
	storageDBCommand.AddCommand(storageDBInitCommand)
}
