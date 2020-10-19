// Copyright © 2020 The Things Industries B.V.

package migrations

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.thethings.network/lorawan-stack/v3/pkg/identityserver/store"
	"go.thethings.network/lorawan-stack/v3/pkg/jsonpb"
	"go.thethings.network/lorawan-stack/v3/pkg/ttipb"
)

// TenantStripeBillingIdentifiers creates the billing identifiers for Stripe based tenants.
type TenantStripeBillingIdentifiers struct{}

func (TenantStripeBillingIdentifiers) Name() string {
	return "tenant_stripe_billing_identifiers"
}

func (TenantStripeBillingIdentifiers) columns() []string {
	return []string{"id", "created_at", "updated_at", "tenant_id", "billing", "billing_identifiers"}
}

func (m TenantStripeBillingIdentifiers) Apply(ctx context.Context, db *gorm.DB) error {
	var models []store.Tenant
	err := db.Model(store.Tenant{}).Select(m.columns()).Find(&models).Error
	if err != nil {
		return err
	}
	for _, model := range models {
		if model.BillingIdentifiers != nil {
			continue
		}
		billing := &ttipb.Billing{}
		if len(model.Billing.RawMessage) > 0 {
			if err := jsonpb.TTN().Unmarshal(model.Billing.RawMessage, billing); err != nil {
				return err
			}
			if billing.GetStripe() == nil {
				continue
			}
		}
		subscriptionID := billing.GetStripe().GetSubscriptionID()
		model.BillingIdentifiers = &subscriptionID
		if err != nil {
			return err
		}
		if err = db.Select(m.columns()).Save(&model).Error; err != nil {
			return err
		}
	}
	return nil
}
func (m TenantStripeBillingIdentifiers) Rollback(ctx context.Context, db *gorm.DB) error {
	// No-op since original identifiers are not deleted during the migration.
	return nil
}
