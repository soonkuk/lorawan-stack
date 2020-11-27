// Copyright © 2020 The Things Industries B.V.

package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/tenant"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"
)

var (
	deleteAuthProviderCommand = &cobra.Command{
		Use:   "delete-auth-provider",
		Short: "Delete a federated authentication provider in the Identity Server database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			logger.Info("Connecting to Identity Server database...")
			db, err := store.Open(ctx, config.IS.DatabaseURI)
			if err != nil {
				return err
			}
			defer db.Close()

			tenantID, err := cmd.Flags().GetString("tenant-id")
			if err != nil {
				return err
			}
			if tenantID == "" {
				tenantID = config.Tenancy.DefaultID
			}
			ctx = tenant.NewContext(ctx, ttipb.TenantIdentifiers{TenantID: tenantID})

			providerID, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			logger.Info("Deleting provider...")
			if err = store.GetAuthenticationProviderStore(db).DeleteAuthenticationProvider(ctx,
				&ttipb.AuthenticationProviderIdentifiers{
					ProviderID: providerID,
				},
			); err != nil {
				return err
			}
			logger.Info("Deleted provider")

			return nil
		},
	}
)

func init() {
	deleteAuthProviderCommand.Flags().String("tenant-id", "", "Tenant ID")
	deleteAuthProviderCommand.Flags().Lookup("tenant-id").Hidden = true
	deleteAuthProviderCommand.Flags().String("id", "", "Provider ID")
	isDBCommand.AddCommand(deleteAuthProviderCommand)
}
