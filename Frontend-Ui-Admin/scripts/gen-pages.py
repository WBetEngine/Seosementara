#!/usr/bin/env python3
"""Generate admin shell pages."""
import os

BASE = os.path.join(os.path.dirname(__file__), "..", "public", "admin")


def shell(title, body, page_title=None):
    pt = page_title or title
    return f"""<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title} — Seosementara Admin</title>
  <link rel="stylesheet" href="/assets/css/admin.css">
  <script src="/assets/js/config.js"></script>
  <script src="https://unpkg.com/htmx.org@2.0.4" defer></script>
  <script src="/assets/js/app.js" defer></script>
</head>
<body>
  <motion id="toast"></motion>
  <div id="drawer-backdrop" hidden></div>
  <aside id="app-drawer" class="app-drawer" aria-hidden="true"></aside>
  <aside id="sidebar" hx-get="/partials/sidebar.html" hx-trigger="load" hx-swap="innerHTML"></aside>
  <div class="admin-layout">
    <div class="admin-main">
      <header class="topbar" hx-get="/partials/topbar.html" hx-trigger="load" hx-swap="innerHTML"></header>
      <main id="main">
        <span id="page-title" class="sr-only">{pt}</span>
        <div class="alert alert--demo">Mode demo — belum ada database, backend, atau Tunnel. Form disimpan ke localStorage.</div>
        {body}
      </main>
    </div>
  </motion>
</body>
</html>""".replace("<motion", "<div").replace("</motion>", "</div>")


def save(rel, html):
    path = os.path.join(BASE, rel)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w", encoding="utf-8") as f:
        f.write(html)


PAGES = {
    "dashboard.html": (
        "Dashboard Admin",
        """
<div class="page-header"><div><h1>Dashboard Admin</h1><p>Ringkasan akun — domain milik & dibagikan</p></div></motion>
<div class="stats-grid">
  <div class="stat-card"><div class="stat-card__label">Domain milik</div><motion class="stat-card__value">12</motion><div class="stat-card__hint">Demo</div></motion>
  <div class="stat-card"><div class="stat-card__label">Dibagikan</div><div class="stat-card__value">3</div></motion>
  <div class="stat-card"><div class="stat-card__label">Undangan</div><div class="stat-card__value">1</div></motion>
  <div class="stat-card"><motion class="stat-card__label">Notifikasi</div><div class="stat-card__value">4</div></motion>
</div>
<div class="card"><div class="card__header">Aktivitas terbaru</div><div class="card__body text-muted">Data muncul setelah backend aktif.</div></motion>
""",
    ),
    "dashboard-domain.html": (
        "Dashboard Domain",
        """
<div class="page-header"><div><h1>Dashboard Domain</h1><p>Per domain portfolio aktif</p></div></motion>
<div class="stats-grid">
  <motion class="stat-card"><div class="stat-card__label">Post</div><div class="stat-card__value">48</div></motion>
  <div class="stat-card"><div class="stat-card__label">Draft</motion><div class="stat-card__value">5</div></motion>
  <div class="stat-card"><div class="stat-card__label">Shortlink</div><div class="stat-card__value">1</div></motion>
  <div class="stat-card"><div class="stat-card__label">Pixel</div><motion class="stat-card__value">—</div></motion>
</div>
""",
    ),
    "dashboard-global.html": (
        "Dashboard Global",
        """
<div class="page-header"><motion><h1>Dashboard Global</h1><p>Super Admin only</p></div></motion>
<div class="stats-grid">
  <div class="stat-card"><div class="stat-card__label">Total domain</div><div class="stat-card__value">3.042</div></motion>
  <div class="stat-card"><div class="stat-card__label">Tunnel</div><div class="stat-card__value"><span class="badge badge--warning">Offline</span></div></motion>
  <div class="stat-card"><div class="stat-card__label">Database</div><div class="stat-card__value"><span class="badge badge--danger">Belum ada</span></motion></motion>
  <div class="stat-card"><div class="stat-card__label">Pekerja</div><div class="stat-card__value">28</div></motion>
</div>
""",
    ),
}

TABLE_DOMAIN = """
<div class="page-header">
  <div><h1>Domain portfolio</h1><p>Kelola ribuan domain native CMS</p></div>
  <div class="page-actions">
    <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/domain-new.html" hx-target="#app-drawer" hx-swap="innerHTML">+ Tambah domain</button>
  </div>
</motion>
<div class="tabs">
  <button class="tab is-active">Domain saya</button>
  <button class="tab">Dibagikan</button>
  <button class="tab">Semua (SA)</button>
</motion>
<div class="card table-wrap">
  <table class="data-table">
    <thead><tr><th>Hostname</th><th>Status</th><th>Owner</th><th></th></tr></thead>
    <tbody>
      <tr>
        <td data-label="Hostname">toko-abc.com</td>
        <td data-label="Status"><span class="badge badge--success">Active</span></td>
        <td data-label="Owner">Anda</td>
        <td data-label="Aksi">
          <button class="btn btn--secondary btn--sm" hx-get="/mock-api/admin/drawers/domain-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button>
        </td>
      </tr>
      <tr>
        <td data-label="Hostname">rezekibelanja.com</td>
        <td data-label="Status"><span class="badge badge--success">Active</span></td>
        <td data-label="Owner">Anda</td>
        <td data-label="Aksi">
          <button class="btn btn--secondary btn--sm" hx-get="/mock-api/admin/drawers/domain-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button>
        </td>
      </tr>
    </tbody>
  </table>
</motion>
"""

PAGES["domain/index.html"] = ("Domain portfolio", TABLE_DOMAIN)

CONTENT_POSTS = """
<div class="page-header">
  <div><h1>Post</h1><p>Konten blog domain aktif</p></div>
  <div class="page-actions">
    <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/post-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">+ Post baru</button>
  </div>
</motion>
<div class="card table-wrap">
  <table class="data-table">
    <thead><tr><th>Judul</th><th>Status</th><th>Diperbarui</th><th></th></tr></thead>
    <tbody>
      <tr>
        <td data-label="Judul">Tips SEO 2026</td>
        <td data-label="Status"><span class="badge badge--success">Published</span></td>
        <td data-label="Diperbarui">21 Mei 2026</td>
        <td data-label="Aksi"><button class="btn btn--secondary btn--sm" hx-get="/mock-api/admin/drawers/post-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></td>
      </tr>
      <tr>
        <td data-label="Judul">Draft artikel</td>
        <td data-label="Status"><span class="badge badge--muted">Draft</span></td>
        <td data-label="Diperbarui">20 Mei 2026</td>
        <td data-label="Aksi"><button class="btn btn--secondary btn--sm" hx-get="/mock-api/admin/drawers/post-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></td>
      </tr>
    </tbody>
  </table>
</motion>
"""

PAGES["content/posts.html"] = ("Post", CONTENT_POSTS)
PAGES["content/pages.html"] = ("Halaman", CONTENT_POSTS.replace("Post", "Halaman").replace("Post baru", "Halaman baru"))
PAGES["content/taxonomy.html"] = (
    "Kategori & tag",
    """
<div class="page-header"><div><h1>Kategori &amp; tag</h1></div>
  <div class="page-actions"><button class="btn btn--primary">+ Kategori</button></div></motion>
<div class="card table-wrap"><table class="data-table"><thead><tr><th>Nama</th><th>Slug</th><th>Post</th></tr></thead>
<tbody><tr><td data-label="Nama">Berita</td><td data-label="Slug">berita</td><td data-label="Post">12</td></tr></tbody></table></motion>
""",
)
PAGES["content/media.html"] = (
    "Media",
    """
<div class="page-header"><div><h1>Perpustakaan media</h1></motion>
  <motion class="page-actions"><button class="btn btn--primary">Upload</button></div></motion>
<div class="card"><div class="card__body empty-state"><p>Belum ada media. Upload setelah backend aktif.</p></div></motion>
""",
)

for name, title in [
    ("seo/meta.html", "Meta & schema"),
    ("seo/sitemap.html", "Sitemap & robots"),
    ("seo/redirects.html", "Redirect manager"),
]:
    PAGES[name] = (
        title,
        f"""
<div class="page-header"><div><h1>{title}</h1><p>Domain portfolio aktif</p></div>
  <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/seo-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></motion>
<div class="card"><div class="card__body"><p class="text-muted">Konfigurasi SEO per domain — isi via drawer.</p></div></motion>
""",
    )

PAGES["plugins/shortlink.html"] = (
    "Shortlink",
    """
<div class="page-header">
  <div><h1>Shortlink</h1><p>url.seosementara.org/{kode}</p></div>
  <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/shortlink-new.html" hx-target="#app-drawer" hx-swap="innerHTML">+ Buat manual</button>
</motion>
<div class="card table-wrap"><table class="data-table"><thead><tr><th>Kode</th><th>Target</th><th>Sumber</th><th>Klik</th></tr></thead>
<tbody>
<tr><td data-label="Kode">rezekibelanja</td><td data-label="Target">https://rezekibelanja.com</td><td data-label="Sumber">Auto</td><td data-label="Klik">128</td></tr>
<tr><td data-label="Kode">promo2025</td><td data-label="Target">https://example.com/p</td><td data-label="Sumber">Manual</td><td data-label="Klik">42</td></tr>
</tbody></table></motion>
""",
)

PAGES["plugins/pixel.html"] = (
    "Pixel Hub",
    """
<div class="page-header"><div><h1>Pixel Hub</h1><p>Kolaborasi Meta, TikTok, Google Ads</p></div></motion>
<div class="stats-grid">
  <div class="stat-card"><div class="stat-card__label">Facebook</div><div class="stat-card__value"><span class="badge badge--muted">Belum setup</span></div></motion>
  <div class="stat-card"><div class="stat-card__label">Antrian</div><motion class="stat-card__value">0</div></motion>
  <div class="stat-card"><div class="stat-card__label">Recovery</div><div class="stat-card__value">—</div></motion>
</div>
<div class="card"><motion class="card__header">Kanal</div><div class="card__body">
  <a class="btn btn--secondary" href="#">Facebook Pro</a>
  <a class="btn btn--secondary" href="#">TikTok</a>
  <a class="btn btn--secondary" href="#">Google Ads</a>
</div></motion>
""",
)

SETTINGS_SUB = """
<div class="settings-layout">
  <nav class="settings-subnav">
    <a href="/admin/settings/index.html">Ringkasan</a>
    <a href="/admin/settings/rbac.html">RBAC</a>
    <a href="/admin/settings/auth.html">Autentikasi</a>
    <a href="/admin/settings/ratelimit.html">Rate limit</a>
    <a href="/admin/settings/ops.html">Operasional</a>
    <a href="/admin/settings/cloudflare.html">Cloudflare</a>
    <a href="/admin/settings/host.html">Host</a>
    <a href="/admin/settings/meta.html">Meta global</a>
    <a href="/admin/settings/notifications.html">Notifikasi</a>
  </nav>
  <div>{body}</div>
</motion>
"""

def settings_page(title, body, active=None):
    sub = SETTINGS_SUB.format(body=body)
    if active:
        sub = sub.replace(f'href="/admin/settings/{active}.html"', f'href="/admin/settings/{active}.html" class="is-active"', 1)
    return title, sub

PAGES["settings/index.html"] = settings_page(
    "Ringkasan sistem",
    """
<div class="page-header"><h1>Ringkasan sistem</h1></motion>
<div class="stats-grid">
  <div class="stat-card"><div class="stat-card__label">Backend</motion><div class="stat-card__value"><span class="badge badge--danger">Offline</span></div></motion>
  <div class="stat-card"><div class="stat-card__label">Database</div><div class="stat-card__value"><span class="badge badge--danger">Belum ada</span></div></motion>
  <motion class="stat-card"><motion class="stat-card__label">Tunnel</div><div class="stat-card__value"><span class="badge badge--warning">Belum connect</span></div></motion>
  <div class="stat-card"><div class="stat-card__label">Versi UI</div><motion class="stat-card__value">0.1</div></motion>
</motion>
""",
    "index",
)

PAGES["settings/rbac.html"] = settings_page(
    "RBAC",
    """
<div class="page-header"><h1>RBAC</h1><button class="btn btn--primary" hx-get="/mock-api/admin/drawers/user-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">+ Pengguna</button></motion>
<div class="card table-wrap"><table class="data-table"><thead><tr><th>Email</th><th>Role</th><th></th></tr></thead>
<tbody><tr><td data-label="Email">admin@seosementara.org</td><td data-label="Role">Super Admin</td>
<td><button class="btn btn--sm btn--secondary" hx-get="/mock-api/admin/drawers/user-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></td></tr></tbody></table></motion>
""",
    "rbac",
)

PAGES["settings/auth.html"] = settings_page(
    "Autentikasi",
    """
<div class="page-header"><h1>Autentikasi &amp; login</h1></motion>
<form class="card card__body form-grid" data-demo-form data-storage-key="sseo_auth_settings">
  <div class="form-group"><label>Panjang password min</label><input class="form-control" name="min_length" type="number" value="10"></motion>
  <div class="form-group"><label>Session TTL (hari)</label><input class="form-control" name="session_days" type="number" value="7"></motion>
  <div class="form-group"><label>Lockout per IP</label><input class="form-control" name="lockout_ip" type="number" value="5"></motion>
  <button type="submit" class="btn btn--primary">Simpan</button>
</form>
""",
    "auth",
)

PAGES["settings/ratelimit.html"] = settings_page(
    "Rate limit",
    """
<div class="page-header"><h1>Rate limit</h1></motion>
<form class="card card__body form-grid form-grid--2" data-demo-form data-storage-key="sseo_ratelimit">
  <div class="form-group"><label>Admin API (rpm/user)</label><input class="form-control" name="admin_rpm" value="300"></motion>
  <div class="form-group"><label>Public API (rpm/IP)</label><input class="form-control" name="public_rpm" value="100"></motion>
  <div class="form-group"><label>Login (rpm/IP)</label><input class="form-control" name="login_rpm" value="5"></motion>
  <button type="submit" class="btn btn--primary">Simpan</button>
</form>
""",
    "ratelimit",
)

PAGES["settings/ops.html"] = settings_page(
    "Operasional",
    """
<div class="page-header"><h1>Operasional</h1></motion>
<form class="card card__body form-grid form-grid--2" data-demo-form data-storage-key="sseo_ops">
  <div class="form-group"><label>Worker concurrency</label><input class="form-control" name="worker_concurrency" value="2"></motion>
  <div class="form-group"><label>Batch size</label><input class="form-control" name="batch_size" value="100"></motion>
  <motion class="form-group"><label>Maintenance mode</label><select class="form-control" name="maintenance"><option value="0">Off</option><option value="1">On</option></select></motion>
  <button type="submit" class="btn btn--primary">Simpan</button>
</form>
""",
    "ops",
)

PAGES["settings/cloudflare.html"] = settings_page(
    "Cloudflare",
    """
<div class="page-header"><h1>Cloudflare</h1>
  <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/cloudflare-token.html" hx-target="#app-drawer" hx-swap="innerHTML">Koneksi API</button></motion>
<div class="card"><div class="card__body">
  <p><strong>Status:</strong> <span class="badge badge--warning">Belum terhubung</span></p>
  <p class="text-muted mt-1">Token disimpan di Workers Secrets — bukan .env di mini PC.</p>
  <hr>
  <h3 class="mb-1">Submenu</h3>
  <ul><li>Domain &amp; env vars</li><li>Tunnel (belum aktif)</li><li>Pages deploy</li><li>DNS</li></ul>
</div></motion>
""",
    "cloudflare",
)

PAGES["settings/host.html"] = settings_page(
    "Host & subdomain",
    """
<div class="page-header"><h1>Host produk</h1>
  <button class="btn btn--primary" hx-get="/mock-api/admin/drawers/host-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">+ Host</button></motion>
<div class="card table-wrap"><table class="data-table"><thead><tr><th>Hostname</th><th>Template</th><th></th></tr></thead>
<tbody>
<tr><td data-label="Hostname">seosementara.org</td><td data-label="Template">apex_default</td><td><button class="btn btn--sm btn--secondary" hx-get="/mock-api/admin/drawers/host-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></td></tr>
<tr><td data-label="Hostname">url.seosementara.org</td><td data-label="Template">subdomain_url</td><td><button class="btn btn--sm btn--secondary" hx-get="/mock-api/admin/drawers/host-edit.html" hx-target="#app-drawer" hx-swap="innerHTML">Edit</button></td></tr>
</tbody></table></motion>
""",
    "host",
)

PAGES["settings/meta.html"] = settings_page(
    "Meta global",
    """
<div class="page-header"><h1>Meta global produk</h1></motion>
<form class="card card__body form-grid" data-demo-form data-storage-key="sseo_meta_global">
  <div class="form-group"><label>Site name</label><input class="form-control" name="site_name" value="Seosementara"></motion>
  <div class="form-group"><label>Default title suffix</label><input class="form-control" name="title_suffix" value="| Seosementara"></motion>
  <div class="form-group"><label>Default description</label><textarea class="form-control" name="description">Platform CMS untuk ribuan domain.</textarea></motion>
  <button type="submit" class="btn btn--primary">Simpan</button>
</form>
""",
    "meta",
)

PAGES["settings/notifications.html"] = settings_page(
    "Notifikasi",
    """
<div class="page-header"><h1>Notifikasi platform</h1></motion>
<form class="card card__body form-grid" data-demo-form data-storage-key="sseo_notifications">
  <div class="form-group"><label>Webhook URL</label><input class="form-control" name="webhook" placeholder="https://..."></motion>
  <div class="form-group"><label>Email alert</label><input class="form-control" name="email" type="email"></motion>
  <button type="submit" class="btn btn--primary">Simpan</button>
</form>
""",
    "notifications",
)

if __name__ == "__main__":
    for rel, (title, body) in PAGES.items():
        save(rel, shell(title, body))
        print("wrote", rel)
    print("done", len(PAGES), "pages")
