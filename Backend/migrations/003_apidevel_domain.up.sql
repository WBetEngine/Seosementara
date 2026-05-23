-- Domain + tunnel defaults untuk apidevel.org (production)

UPDATE domain_env_config SET value = 'apidevel.org', updated_at = now()
WHERE key = 'PRIMARY_DOMAIN';

UPDATE domain_env_config SET value = 'https://apidevel.org', updated_at = now()
WHERE key = 'APEX_URL';

UPDATE domain_env_config SET value = 'https://api.apidevel.org', updated_at = now()
WHERE key = 'API_BASE_URL';

UPDATE cloudflare_tunnel_config SET
    tunnel_id = '9c9882ac-d826-4714-92b0-6531e6a9844a',
    tunnel_name = 'seosementara-api',
    is_active = true,
    origin_url = 'http://localhost:8080',
    connector_status = 'healthy',
    updated_at = now()
WHERE id = 1;

INSERT INTO cloudflare_tunnel_routes (hostname, path_prefix, service_url, is_enabled)
VALUES ('api.apidevel.org', NULL, 'http://localhost:8080', true)
ON CONFLICT (hostname, path_prefix) DO UPDATE SET
    service_url = EXCLUDED.service_url,
    is_enabled = true;

UPDATE cloudflare_pages_projects SET
    pages_url = 'https://seosementara.seosementara3.workers.dev',
    updated_at = now()
WHERE project_type = 'admin';
