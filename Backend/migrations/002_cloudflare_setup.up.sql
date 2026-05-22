-- Cloudflare Settings (Plan/15) — singleton credentials, tunnel, pages, env, audit

CREATE TABLE IF NOT EXISTS cloudflare_credentials (
    id                 SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    auth_type          TEXT NOT NULL CHECK (auth_type IN ('api_token', 'global_api_key')),
    token_ciphertext   BYTEA NOT NULL,
    token_nonce        BYTEA NOT NULL,
    account_id         TEXT,
    account_email      TEXT,
    last_validated_at  TIMESTAMPTZ,
    validation_error   TEXT,
    updated_by         BIGINT,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS domain_env_config (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    is_public   BOOLEAN NOT NULL DEFAULT true,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO domain_env_config (key, value, is_public) VALUES
    ('PRIMARY_DOMAIN', 'seosementara.org', true),
    ('APEX_URL', 'https://seosementara.org', true),
    ('API_BASE_URL', 'https://seosementara.org', true),
    ('ADMIN_BASE_PATH', '/admin', true),
    ('ENVIRONMENT', 'production', true),
    ('TURNSTILE_SITE_KEY', '', true),
    ('CDN_ASSETS_URL', '', true),
    ('ZONE_ID', '', true)
ON CONFLICT (key) DO NOTHING;

CREATE TABLE IF NOT EXISTS cloudflare_tunnel_config (
    id                 SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    tunnel_id          TEXT,
    tunnel_name        TEXT NOT NULL DEFAULT 'seosementara-api',
    install_command    TEXT,
    is_active          BOOLEAN NOT NULL DEFAULT false,
    connector_status   TEXT NOT NULL DEFAULT 'unknown',
    last_seen_at       TIMESTAMPTZ,
    origin_url         TEXT NOT NULL DEFAULT 'http://127.0.0.1:8080',
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO cloudflare_tunnel_config (id, tunnel_name, origin_url)
VALUES (1, 'seosementara-api', 'http://127.0.0.1:8080')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS cloudflare_tunnel_routes (
    id           BIGSERIAL PRIMARY KEY,
    hostname     TEXT NOT NULL,
    path_prefix  TEXT,
    service_url  TEXT NOT NULL,
    is_enabled   BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (hostname, path_prefix)
);

INSERT INTO cloudflare_tunnel_routes (hostname, path_prefix, service_url) VALUES
    ('seosementara.org', '/api', 'http://127.0.0.1:8080'),
    ('seosementara.org', '/admin', 'http://127.0.0.1:8080'),
    ('*.seosementara.org', NULL, 'http://127.0.0.1:8080')
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS cloudflare_pages_projects (
    id                 BIGSERIAL PRIMARY KEY,
    project_type       TEXT NOT NULL CHECK (project_type IN ('admin', 'public')),
    project_id         TEXT,
    project_name       TEXT NOT NULL,
    production_branch  TEXT NOT NULL DEFAULT 'main',
    root_directory     TEXT,
    build_command      TEXT,
    deploy_command     TEXT,
    pages_url          TEXT,
    custom_domains     JSONB NOT NULL DEFAULT '[]',
    env_synced_at      TIMESTAMPTZ,
    last_deploy_at     TIMESTAMPTZ,
    deploy_status      TEXT,
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_type)
);

INSERT INTO cloudflare_pages_projects (
    project_type, project_name, production_branch, root_directory, build_command, deploy_command, pages_url
) VALUES (
    'admin', 'seosementara', 'main', 'Frontend-admin', 'npm ci', 'npx wrangler deploy',
    'https://seosementara.seosementara3.workers.dev'
) ON CONFLICT (project_type) DO NOTHING;

CREATE TABLE IF NOT EXISTS cloudflare_api_logs (
    id            BIGSERIAL PRIMARY KEY,
    method        TEXT NOT NULL,
    path          TEXT NOT NULL,
    status_code   INT,
    success       BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT,
    duration_ms   INT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cloudflare_api_logs_created
    ON cloudflare_api_logs (created_at DESC);
