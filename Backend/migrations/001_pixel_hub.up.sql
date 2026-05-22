-- Pixel Hub foundation + Facebook Pro tables

CREATE TABLE IF NOT EXISTS pixel_hub_settings (
    id                 BIGSERIAL PRIMARY KEY,
    tracking_hostname  TEXT NOT NULL DEFAULT 'pelacak.seosementara.org',
    default_mode       TEXT NOT NULL DEFAULT 'server_first'
        CHECK (default_mode IN ('server_first', 'hybrid', 'legacy_client')),
    consent_required   BOOLEAN NOT NULL DEFAULT false,
    collect_path       TEXT NOT NULL DEFAULT '/collect',
    script_version     TEXT NOT NULL DEFAULT '1',
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO pixel_hub_settings (id) VALUES (1) ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS pixel_credentials (
    id                  BIGSERIAL PRIMARY KEY,
    platform            TEXT NOT NULL DEFAULT 'facebook',
    name                TEXT NOT NULL,
    secret_ciphertext   BYTEA NOT NULL,
    secret_nonce        BYTEA NOT NULL,
    last_validated_at   TIMESTAMPTZ,
    validation_status   TEXT,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pixel_configs (
    id                  BIGSERIAL PRIMARY KEY,
    platform            TEXT NOT NULL DEFAULT 'facebook'
        CHECK (platform IN ('facebook', 'tiktok', 'gads')),
    scope               TEXT NOT NULL DEFAULT 'global'
        CHECK (scope IN ('global', 'managed_domain', 'shortlink')),
    managed_domain_id   BIGINT,
    name                TEXT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    mode_override       TEXT CHECK (mode_override IN ('server_first', 'hybrid', 'legacy_client')),
    external_ids        JSONB NOT NULL DEFAULT '{}',
    credentials_id      BIGINT REFERENCES pixel_credentials(id) ON DELETE SET NULL,
    capi_enabled        BOOLEAN NOT NULL DEFAULT true,
    browser_pixel_enabled BOOLEAN NOT NULL DEFAULT false,
    test_event_code     TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pixel_configs_platform_active
    ON pixel_configs (platform, is_active);

CREATE TABLE IF NOT EXISTS pixel_domain_assignments (
    id                  BIGSERIAL PRIMARY KEY,
    pixel_config_id     BIGINT NOT NULL REFERENCES pixel_configs(id) ON DELETE CASCADE,
    managed_domain_id   BIGINT NOT NULL,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    deployed_at         TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (pixel_config_id, managed_domain_id)
);

CREATE TABLE IF NOT EXISTS pixel_events (
    id                  BIGSERIAL PRIMARY KEY,
    canonical_event     TEXT,
    platform            TEXT NOT NULL DEFAULT 'facebook',
    pixel_config_id     BIGINT REFERENCES pixel_configs(id) ON DELETE SET NULL,
    event_name          TEXT NOT NULL,
    event_id            TEXT NOT NULL,
    managed_domain_id   BIGINT,
    url_link_id         BIGINT,
    payload             JSONB NOT NULL DEFAULT '{}',
    status              TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'sent', 'failed', 'dropped_bot', 'skipped')),
    platform_event_id   TEXT,
    error_message       TEXT,
    sent_at             TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pixel_events_pending
    ON pixel_events (created_at)
    WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS idx_pixel_events_facebook_diag
    ON pixel_events (platform, status, created_at DESC);

CREATE TABLE IF NOT EXISTS pixel_facebook_stats_daily (
    stat_date           DATE NOT NULL,
    pixel_config_id     BIGINT REFERENCES pixel_configs(id) ON DELETE CASCADE,
    events_received     INT NOT NULL DEFAULT 0,
    events_sent         INT NOT NULL DEFAULT 0,
    events_failed       INT NOT NULL DEFAULT 0,
    PRIMARY KEY (stat_date, pixel_config_id)
);
