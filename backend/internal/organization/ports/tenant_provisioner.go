package ports

import "context"

type TenantProvisioner interface {
	CreateSchema(ctx context.Context, schemaName string) error
}
