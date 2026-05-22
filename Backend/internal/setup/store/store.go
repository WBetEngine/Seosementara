package store

import (
	"context"
	"time"
)

type SetupStore interface {
	GetCredentials(ctx context.Context) (*Credentials, error)
	GetCredentialsSecret(ctx context.Context, encKey []byte) (authType, secret string, err error)
	SaveCredentials(ctx context.Context, authType string, secret []byte, accountID, accountEmail string, encKey []byte) error
	UpdateCredentialsValidation(ctx context.Context, accountID string, ok bool, errMsg string) error

	ListEnv(ctx context.Context) ([]EnvEntry, error)
	UpsertEnv(ctx context.Context, key, value string, isPublic bool) error
	GetEnvMap(ctx context.Context) (map[string]string, error)

	GetTunnelConfig(ctx context.Context) (*TunnelConfig, error)
	SaveTunnelConfig(ctx context.Context, tunnelID, name, installCmd, origin string, active bool, status string) error
	UpdateTunnelStatus(ctx context.Context, status string, lastSeen *time.Time) error
	ListTunnelRoutes(ctx context.Context) ([]TunnelRoute, error)
	UpsertTunnelRoute(ctx context.Context, r TunnelRoute) (int64, error)
	DeleteTunnelRoute(ctx context.Context, id int64) error

	GetPagesProject(ctx context.Context, projectType string) (*PagesProject, error)
	SavePagesProject(ctx context.Context, p PagesProject) error
	UpdatePagesDeploy(ctx context.Context, projectType, status string, deployedAt time.Time) error

	InsertAPILog(ctx context.Context, method, path string, status int, success bool, errMsg string, durationMs int) error
	ListAPILogs(ctx context.Context, limit int) ([]APILog, error)
}
