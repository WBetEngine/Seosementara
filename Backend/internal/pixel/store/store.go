package store

import "context"

// PixelStore abstracts pixel hub persistence.
type PixelStore interface {
	GetHubSettings(ctx context.Context) (*HubSettings, error)
	UpdateHubSettings(ctx context.Context, hostname, mode string, consent bool) error

	GetFacebookConfig(ctx context.Context) (*FacebookConfig, error)
	SaveFacebookSetup(ctx context.Context, in FacebookSetupInput, encKey []byte) (*FacebookConfig, error)
	GetCredentialSecret(ctx context.Context, credID int64, encKey []byte) (string, error)
	UpdateCredentialValidation(ctx context.Context, credID int64, status string) error

	EnqueueEvent(ctx context.Context, ev PixelEvent) (int64, error)
	ListPendingEvents(ctx context.Context, limit int) ([]PixelEvent, error)
	MarkEventSent(ctx context.Context, id int64, platformEventID string) error
	MarkEventFailed(ctx context.Context, id int64, errMsg string) error
	ListEvents(ctx context.Context, status string, limit, offset int) ([]PixelEvent, int64, error)
	GetDiagnostics(ctx context.Context) (*Diagnostics, error)
	IncrementDailyStat(ctx context.Context, configID int64, field string) error

	ListDomainAssignments(ctx context.Context, configID int64) ([]DomainAssignment, error)
	AssignDomain(ctx context.Context, configID, domainID int64, hostname string) error
	UnassignDomain(ctx context.Context, configID, domainID int64) error
}
