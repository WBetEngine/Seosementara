package store

import (
	"context"
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

func (p *PostgresStore) GetCredentials(ctx context.Context) (*Credentials, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT id, auth_type, account_id, COALESCE(account_email,''),
			last_validated_at, COALESCE(validation_error,''), updated_at
		FROM cloudflare_credentials WHERE id = 1`)
	var c Credentials
	err := row.Scan(&c.ID, &c.AuthType, &c.AccountID, &c.AccountEmail,
		&c.LastValidatedAt, &c.ValidationError, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.TokenMasked = "••••••••"
	return &c, nil
}

func (p *PostgresStore) GetCredentialsSecret(ctx context.Context, encKey []byte) (authType, secret string, err error) {
	row := p.pool.QueryRow(ctx, `
		SELECT auth_type, token_ciphertext, token_nonce FROM cloudflare_credentials WHERE id = 1`)
	var ct, nonce []byte
	err = row.Scan(&authType, &ct, &nonce)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", err
	}
	plain, err := crypto.DecryptAESGCM(encKey, ct, nonce)
	return authType, string(plain), err
}

func (p *PostgresStore) SaveCredentials(ctx context.Context, authType string, secret []byte, accountID, accountEmail string, encKey []byte) error {
	ct, nonce, err := crypto.EncryptAESGCM(encKey, secret)
	if err != nil {
		return err
	}
	_, err = p.pool.Exec(ctx, `
		INSERT INTO cloudflare_credentials (id, auth_type, token_ciphertext, token_nonce, account_id, account_email, updated_at)
		VALUES (1, $1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO UPDATE SET
			auth_type = EXCLUDED.auth_type,
			token_ciphertext = EXCLUDED.token_ciphertext,
			token_nonce = EXCLUDED.token_nonce,
			account_id = EXCLUDED.account_id,
			account_email = EXCLUDED.account_email,
			validation_error = NULL,
			updated_at = now()`,
		authType, ct, nonce, nullStr(accountID), nullStr(accountEmail))
	return err
}

func (p *PostgresStore) UpdateCredentialsValidation(ctx context.Context, accountID string, ok bool, errMsg string) error {
	var validated *time.Time
	if ok {
		t := time.Now()
		validated = &t
	}
	_, err := p.pool.Exec(ctx, `
		UPDATE cloudflare_credentials SET
			account_id = COALESCE(NULLIF($1,''), account_id),
			last_validated_at = $2,
			validation_error = $3,
			updated_at = now()
		WHERE id = 1`, accountID, validated, nullStr(errMsg))
	return err
}

func (p *PostgresStore) ListEnv(ctx context.Context) ([]EnvEntry, error) {
	rows, err := p.pool.Query(ctx, `SELECT key, value, is_public, updated_at FROM domain_env_config ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EnvEntry
	for rows.Next() {
		var e EnvEntry
		if err := rows.Scan(&e.Key, &e.Value, &e.IsPublic, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (p *PostgresStore) UpsertEnv(ctx context.Context, key, value string, isPublic bool) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO domain_env_config (key, value, is_public, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, is_public = EXCLUDED.is_public, updated_at = now()`,
		key, value, isPublic)
	return err
}

func (p *PostgresStore) GetEnvMap(ctx context.Context) (map[string]string, error) {
	rows, err := p.pool.Query(ctx, `SELECT key, value FROM domain_env_config`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		m[k] = v
	}
	return m, rows.Err()
}

func (p *PostgresStore) GetTunnelConfig(ctx context.Context) (*TunnelConfig, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT id, COALESCE(tunnel_id,''), tunnel_name, COALESCE(install_command,''),
			is_active, connector_status, last_seen_at, origin_url, updated_at
		FROM cloudflare_tunnel_config WHERE id = 1`)
	var t TunnelConfig
	err := row.Scan(&t.ID, &t.TunnelID, &t.TunnelName, &t.InstallCommand,
		&t.IsActive, &t.ConnectorStatus, &t.LastSeenAt, &t.OriginURL, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &t, err
}

func (p *PostgresStore) SaveTunnelConfig(ctx context.Context, tunnelID, name, installCmd, origin string, active bool, status string) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE cloudflare_tunnel_config SET
			tunnel_id = $1, tunnel_name = $2, install_command = $3, origin_url = $4,
			is_active = $5, connector_status = $6, updated_at = now()
		WHERE id = 1`,
		nullStr(tunnelID), name, installCmd, origin, active, status)
	return err
}

func (p *PostgresStore) UpdateTunnelStatus(ctx context.Context, status string, lastSeen *time.Time) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE cloudflare_tunnel_config SET connector_status = $1, last_seen_at = $2, updated_at = now() WHERE id = 1`,
		status, lastSeen)
	return err
}

func (p *PostgresStore) ListTunnelRoutes(ctx context.Context) ([]TunnelRoute, error) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, hostname, COALESCE(path_prefix,''), service_url, is_enabled
		FROM cloudflare_tunnel_routes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TunnelRoute
	for rows.Next() {
		var r TunnelRoute
		var path string
		if err := rows.Scan(&r.ID, &r.Hostname, &path, &r.ServiceURL, &r.IsEnabled); err != nil {
			return nil, err
		}
		r.PathPrefix = path
		out = append(out, r)
	}
	return out, rows.Err()
}

func (p *PostgresStore) UpsertTunnelRoute(ctx context.Context, r TunnelRoute) (int64, error) {
	if r.ID > 0 {
		_, err := p.pool.Exec(ctx, `
			UPDATE cloudflare_tunnel_routes SET hostname=$1, path_prefix=NULLIF($2,''), service_url=$3, is_enabled=$4
			WHERE id=$5`, r.Hostname, r.PathPrefix, r.ServiceURL, r.IsEnabled, r.ID)
		return r.ID, err
	}
	var id int64
	err := p.pool.QueryRow(ctx, `
		INSERT INTO cloudflare_tunnel_routes (hostname, path_prefix, service_url, is_enabled)
		VALUES ($1, NULLIF($2,''), $3, $4)
		ON CONFLICT (hostname, path_prefix) DO UPDATE SET service_url = EXCLUDED.service_url, is_enabled = EXCLUDED.is_enabled
		RETURNING id`, r.Hostname, r.PathPrefix, r.ServiceURL, r.IsEnabled).Scan(&id)
	return id, err
}

func (p *PostgresStore) DeleteTunnelRoute(ctx context.Context, id int64) error {
	_, err := p.pool.Exec(ctx, `DELETE FROM cloudflare_tunnel_routes WHERE id = $1`, id)
	return err
}

func (p *PostgresStore) GetPagesProject(ctx context.Context, projectType string) (*PagesProject, error) {
	row := p.pool.QueryRow(ctx, `
		SELECT id, project_type, COALESCE(project_id,''), project_name, production_branch,
			COALESCE(root_directory,''), COALESCE(build_command,''), COALESCE(deploy_command,''),
			COALESCE(pages_url,''), COALESCE(deploy_status,''), last_deploy_at, env_synced_at, updated_at
		FROM cloudflare_pages_projects WHERE project_type = $1`, projectType)
	var pp PagesProject
	err := row.Scan(&pp.ID, &pp.ProjectType, &pp.ProjectID, &pp.ProjectName, &pp.ProductionBranch,
		&pp.RootDirectory, &pp.BuildCommand, &pp.DeployCommand, &pp.PagesURL, &pp.DeployStatus,
		&pp.LastDeployAt, &pp.EnvSyncedAt, &pp.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &pp, err
}

func (p *PostgresStore) SavePagesProject(ctx context.Context, pp PagesProject) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO cloudflare_pages_projects (
			project_type, project_id, project_name, production_branch, root_directory,
			build_command, deploy_command, pages_url, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())
		ON CONFLICT (project_type) DO UPDATE SET
			project_id = EXCLUDED.project_id,
			project_name = EXCLUDED.project_name,
			production_branch = EXCLUDED.production_branch,
			root_directory = EXCLUDED.root_directory,
			build_command = EXCLUDED.build_command,
			deploy_command = EXCLUDED.deploy_command,
			pages_url = EXCLUDED.pages_url,
			updated_at = now()`,
		pp.ProjectType, nullStr(pp.ProjectID), pp.ProjectName, pp.ProductionBranch,
		pp.RootDirectory, pp.BuildCommand, pp.DeployCommand, pp.PagesURL)
	return err
}

func (p *PostgresStore) UpdatePagesDeploy(ctx context.Context, projectType, status string, deployedAt time.Time) error {
	_, err := p.pool.Exec(ctx, `
		UPDATE cloudflare_pages_projects SET deploy_status = $1, last_deploy_at = $2, updated_at = now()
		WHERE project_type = $3`, status, deployedAt, projectType)
	return err
}

func (p *PostgresStore) InsertAPILog(ctx context.Context, method, path string, status int, success bool, errMsg string, durationMs int) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO cloudflare_api_logs (method, path, status_code, success, error_message, duration_ms)
		VALUES ($1,$2,$3,$4,$5,$6)`, method, path, status, success, nullStr(errMsg), durationMs)
	return err
}

func (p *PostgresStore) ListAPILogs(ctx context.Context, limit int) ([]APILog, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := p.pool.Query(ctx, `
		SELECT id, method, path, COALESCE(status_code,0), success, COALESCE(error_message,''), COALESCE(duration_ms,0), created_at
		FROM cloudflare_api_logs ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []APILog
	for rows.Next() {
		var l APILog
		if err := rows.Scan(&l.ID, &l.Method, &l.Path, &l.StatusCode, &l.Success, &l.ErrorMessage, &l.DurationMs, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
