package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (p *PostgresStore) GetHubSettings(ctx context.Context) (*HubSettings, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT id, tracking_hostname, default_mode, consent_required, collect_path, script_version
		FROM pixel_hub_settings WHERE id = 1`)
	var h HubSettings
	err := row.Scan(&h.ID, &h.TrackingHostname, &h.DefaultMode, &h.ConsentRequired, &h.CollectPath, &h.ScriptVersion)
	return &h, err
}

func (p *PostgresStore) UpdateHubSettings(ctx context.Context, hostname, mode string, consent bool) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE pixel_hub_settings SET tracking_hostname=$1, default_mode=$2, consent_required=$3, updated_at=now()
		WHERE id=1`, hostname, mode, consent)
	return err
}

func (p *PostgresStore) GetFacebookConfig(ctx context.Context) (*FacebookConfig, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT c.id, c.name, c.scope, c.managed_domain_id, c.is_active, c.mode_override,
			c.external_ids->>'pixel_id', COALESCE(c.external_ids->>'business_id',''),
			c.capi_enabled, c.browser_pixel_enabled, COALESCE(c.test_event_code,''),
			c.credentials_id, COALESCE(cr.name,''), COALESCE(cr.validation_status,'unknown'),
			cr.last_validated_at, c.updated_at
		FROM pixel_configs c
		LEFT JOIN pixel_credentials cr ON cr.id = c.credentials_id
		WHERE c.platform = 'facebook' AND c.scope = 'global'
		ORDER BY c.id LIMIT 1`)
	var c FacebookConfig
	var mode *string
	var credID *int64
	var lastVal *time.Time
	err := row.Scan(&c.ID, &c.Name, &c.Scope, &c.ManagedDomainID, &c.IsActive, &mode,
		&c.PixelID, &c.BusinessID, &c.CAPIEnabled, &c.BrowserPixelEnabled, &c.TestEventCode,
		&credID, &c.CredentialName, &c.ValidationStatus, &lastVal, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	c.ModeOverride = mode
	c.CredentialsID = credID
	c.LastValidatedAt = lastVal
	return &c, err
}

func (p *PostgresStore) SaveFacebookSetup(ctx context.Context, in FacebookSetupInput, encKey []byte) (*FacebookConfig, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var credID *int64
	if in.CAPIAccessToken != "" {
		ct, nonce, err := crypto.EncryptAESGCM(encKey, []byte(in.CAPIAccessToken))
		if err != nil {
			return nil, err
		}
		name := in.CredentialName
		if name == "" {
			name = "Meta CAPI default"
		}
		err = tx.QueryRow(ctx, `
			INSERT INTO pixel_credentials (platform, name, secret_ciphertext, secret_nonce, validation_status)
			VALUES ('facebook', $1, $2, $3, 'unknown') RETURNING id`,
			name, ct, nonce).Scan(&credID)
		if err != nil {
			return nil, err
		}
	}

	ext, _ := json.Marshal(map[string]string{
		"pixel_id": in.PixelID, "business_id": in.BusinessID,
	})
	var exists bool
	_ = tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM pixel_configs WHERE platform='facebook' AND scope='global')`).
		Scan(&exists)
	if exists {
		_, err = tx.Exec(ctx, `
			UPDATE pixel_configs SET name=$1, external_ids=$2,
				credentials_id=COALESCE($3, credentials_id),
				capi_enabled=$4, browser_pixel_enabled=$5, test_event_code=$6, updated_at=now()
			WHERE platform='facebook' AND scope='global'`,
			in.Name, ext, credID, in.CAPIEnabled, in.BrowserPixelEnabled, in.TestEventCode)
	} else {
		_, err = tx.Exec(ctx, `
			INSERT INTO pixel_configs (platform, scope, name, is_active, external_ids, credentials_id,
				capi_enabled, browser_pixel_enabled, test_event_code)
			VALUES ('facebook', $1, $2, true, $3, $4, $5, $6, $7)`,
			in.Scope, in.Name, ext, credID, in.CAPIEnabled, in.BrowserPixelEnabled, in.TestEventCode)
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return p.GetFacebookConfig(ctx)
}

func (p *PostgresStore) GetCredentialSecret(ctx context.Context, credID int64, encKey []byte) (string, error) {
	var ct, nonce []byte
	err := p.pool.QueryRow(ctx, `
		SELECT secret_ciphertext, secret_nonce FROM pixel_credentials WHERE id=$1`, credID).
		Scan(&ct, &nonce)
	if err != nil {
		return "", err
	}
	plain, err := crypto.DecryptAESGCM(encKey, ct, nonce)
	return string(plain), err
}

func (p *PostgresStore) UpdateCredentialValidation(ctx context.Context, credID int64, status string) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE pixel_credentials SET validation_status=$1, last_validated_at=now() WHERE id=$2`,
		status, credID)
	return err
}

func (p *PostgresStore) EnqueueEvent(ctx context.Context, ev PixelEvent) (int64, error) {
	var id int64
	err := p.pool.QueryRow(ctx, `
		INSERT INTO pixel_events (platform, pixel_config_id, event_name, event_id, managed_domain_id, payload, status)
		VALUES ('facebook', $1, $2, $3, $4, $5, 'pending') RETURNING id`,
		ev.PixelConfigID, ev.EventName, ev.EventID, ev.ManagedDomainID, ev.Payload).Scan(&id)
	return id, err
}

func (p *PostgresStore) ListPendingEvents(ctx context.Context, limit int) ([]PixelEvent, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, event_name, event_id, status, pixel_config_id, managed_domain_id,
			COALESCE(error_message,''), COALESCE(platform_event_id,''), created_at, payload
		FROM pixel_events WHERE status='pending' ORDER BY created_at LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

func (p *PostgresStore) MarkEventSent(ctx context.Context, id int64, platformEventID string) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE pixel_events SET status='sent', platform_event_id=$1, sent_at=now() WHERE id=$2`,
		platformEventID, id)
	return err
}

func (p *PostgresStore) MarkEventFailed(ctx context.Context, id int64, errMsg string) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE pixel_events SET status='failed', error_message=$1 WHERE id=$2`, errMsg, id)
	return err
}

func (p *PostgresStore) ListEvents(ctx context.Context, status string, limit, offset int) ([]PixelEvent, int64, error) {
	var total int64
	qCount := `SELECT COUNT(*) FROM pixel_events WHERE platform='facebook'`
	qList := `SELECT id, event_name, event_id, status, pixel_config_id, managed_domain_id,
		COALESCE(error_message,''), COALESCE(platform_event_id,''), created_at, payload
		FROM pixel_events WHERE platform='facebook'`
	args := []any{}
	if status != "" {
		qCount += ` AND status=$1`
		qList += ` AND status=$1`
		args = append(args, status)
	}
	if err := p.pool.QueryRow(ctx, qCount, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	qList += ` ORDER BY created_at DESC LIMIT $` + itoa(len(args)+1) + ` OFFSET $` + itoa(len(args)+2)
	args = append(args, limit, offset)
	rows, err := p.pool.Query(ctx, qList, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	evs, err := scanEvents(rows)
	return evs, total, err
}

func (p *PostgresStore) GetDiagnostics(ctx context.Context) (*Diagnostics, error) {
	d := &Diagnostics{ConnectionState: "unknown"}
	_ = p.pool.QueryRow(ctx, `
		SELECT COALESCE(validation_status,'unknown') FROM pixel_credentials
		WHERE platform='facebook' ORDER BY updated_at DESC LIMIT 1`).Scan(&d.ConnectionState)

	_ = p.pool.QueryRow(ctx, `
		SELECT COUNT(*) FILTER (WHERE status='pending'),
			COUNT(*) FILTER (WHERE status='failed' AND created_at > now()-interval '24 hours'),
			COUNT(*) FILTER (WHERE status='sent' AND created_at > now()-interval '24 hours'),
			COUNT(*) FILTER (WHERE created_at > now()-interval '24 hours')
		FROM pixel_events WHERE platform='facebook'`).
		Scan(&d.PendingCount, &d.Failed24h, &d.Sent24h, &d.Received24h)

	if d.Received24h > 0 {
		d.FailureRatePct = float64(d.Failed24h) / float64(d.Received24h) * 100
	}
	_ = p.pool.QueryRow(ctx, `
		SELECT COALESCE(error_message,'') FROM pixel_events
		WHERE platform='facebook' AND status='failed' ORDER BY created_at DESC LIMIT 1`).
		Scan(&d.LastError)
	return d, nil
}

func (p *PostgresStore) IncrementDailyStat(ctx context.Context, configID int64, field string) error {
	col := "events_received"
	switch field {
	case "sent":
		col = "events_sent"
	case "failed":
		col = "events_failed"
	}
	_, err := p.pool.Exec(ctx, `
		INSERT INTO pixel_facebook_stats_daily (stat_date, pixel_config_id, `+col+`)
		VALUES (CURRENT_DATE, $1, 1)
		ON CONFLICT (stat_date, pixel_config_id) DO UPDATE SET `+col+` = pixel_facebook_stats_daily.`+col+` + 1`,
		configID)
	return err
}

func (p *PostgresStore) ListDomainAssignments(ctx context.Context, configID int64) ([]DomainAssignment, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT a.id, a.pixel_config_id, a.managed_domain_id,
			COALESCE(a.managed_domain_id::text, 'domain-'||a.managed_domain_id::text),
			a.is_active, a.deployed_at
		FROM pixel_domain_assignments a WHERE a.pixel_config_id=$1`, configID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DomainAssignment
	for rows.Next() {
		var a DomainAssignment
		if err := rows.Scan(&a.ID, &a.PixelConfigID, &a.ManagedDomainID, &a.DomainHostname, &a.IsActive, &a.DeployedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (p *PostgresStore) AssignDomain(ctx context.Context, configID, domainID int64, hostname string) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO pixel_domain_assignments (pixel_config_id, managed_domain_id)
		VALUES ($1, $2) ON CONFLICT (pixel_config_id, managed_domain_id) DO UPDATE SET is_active=true`,
		configID, domainID)
	return err
}

func (p *PostgresStore) UnassignDomain(ctx context.Context, configID, domainID int64) error {
	_, err := p.pool.Exec(ctx, `
		DELETE FROM pixel_domain_assignments WHERE pixel_config_id=$1 AND managed_domain_id=$2`,
		configID, domainID)
	return err
}

func itoa(n int) string {
	if n == 1 {
		return "1"
	}
	if n == 2 {
		return "2"
	}
	return "3"
}

type rowScanner interface {
	Next() bool
	Scan(dest ...any) error
	Close()
}

func scanEvents(rows rowScanner) ([]PixelEvent, error) {
	var out []PixelEvent
	for rows.Next() {
		var e PixelEvent
		if err := rows.Scan(&e.ID, &e.EventName, &e.EventID, &e.Status, &e.PixelConfigID,
			&e.ManagedDomainID, &e.ErrorMessage, &e.PlatformEventID, &e.CreatedAt, &e.Payload); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

var _ PixelStore = (*PostgresStore)(nil)
