# 11 — RBAC & Permission Share Domain

> Melengkapi matriks singkat di [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md).  
> Database: [10-database-postgresql.md](./10-database-postgresql.md)

## 1. Dua Lapisan Hak Akses

```mermaid
flowchart TB
  subgraph system [Lapisan 1 - Sistem]
    SA[super_admin]
    W[worker]
  end
  subgraph domain [Lapisan 2 - Per Domain Portfolio]
    O[owner]
    S[shared user + permission checklist]
  end
  SA --> Global
  W --> O
  W --> S
  O --> Full pada domain milik
  S --> Sesuai checklist
```

| Lapisan | Scope | Ditentukan oleh |
|---------|-------|-----------------|
| **Sistem** | Seluruh platform | `users.role` — hanya `super_admin` \| `worker` |
| **Domain** | Satu `managed_domain` | Ownership + `domain_shares.permissions` |

**Prinsip:** RBAC domain **tidak** menggantikan ownership; share hanya menambah hak terbatas pada domain orang lain.

---

## 2. Peran Sistem (Lapisan 1)

| Role | Kode | Kemampuan global |
|------|------|------------------|
| **Super Admin** | `super_admin` | Semua domain, Setup Host/subdomain, Setup Backend, kelola user sistem, transfer ownership |
| **Worker** | `worker` | Hanya domain milik + yang di-share; tidak ada menu Setup sistem |

Tidak ada peran sistem `site_manager`, `editor`, dll. — itu diganti **permission checklist per share** (lapisan 2).

---

## 3. Share Domain — Mode & Preset

Saat owner / co-admin (setelah disetujui) mengundang admin lain:

### 3.1 Mode cepat (preset)

| Preset | Label UI | Ringkasan |
|--------|----------|-----------|
| `read_only` | Hanya baca | Lihat domain, konten, media, SEO, laporan — **tanpa** ubah apapun |
| `edit` | Bisa edit | Tulis/edit konten, media, SEO — **tanpa** publish, hapus, bulk, share, setting domain |
| `full_edit` | Edit + publish | Edit + publish/unpublish + media penuh — **tanpa** hapus permanen, bulk, share, setting |
| `co_admin` | Co-Admin | Hampir penuh pada domain ini; boleh undang orang lain (pending owner) |
| `custom` | Kustom | Centang manual per permission di checklist §4 |

### 3.2 Alur UI (`/admin/sites/{id}/sharing`)

```
[ ] Preset: ( ) Read only  ( ) Edit  ( ) Full edit  ( ) Co-Admin  (•) Kustom

Checklist permission:
  [x] Lihat domain & dashboard domain
  [x] Lihat daftar post / page
  [ ] Buat / edit post / page
  ...

[Simpan]  → Owner: aktif langsung | Co-admin: pending approval
```

---

## 4. Checklist Permission (Granular)

Setiap permission = satu key boolean. Disimpan sebagai `JSONB` di `domain_shares.permissions`.

### 4.1 Domain & navigasi

| Key | Label | Read only | Edit | Co-Admin |
|-----|-------|:---------:|:----:|:--------:|
| `domain.view` | Lihat info domain & masuk site switcher | ✓ | ✓ | ✓ |
| `domain.settings.view` | Lihat pengaturan domain | ✓ | ✓ | ✓ |
| `domain.settings.edit` | Ubah pengaturan domain | | | ✓ |
| `domain.share` | Undang / kelola share user lain | | | ✓* |

\* Co-admin: undangan co-admin tetap **pending** sampai owner setujui ([09](./09-model-domain-host-dan-subdomain.md)).

### 4.2 Konten

| Key | Label | Read only | Edit | Full edit | Co-Admin |
|-----|-------|:---------:|:----:|:---------:|:--------:|
| `content.posts.view` | Lihat post | ✓ | ✓ | ✓ | ✓ |
| `content.posts.create` | Buat post | | ✓ | ✓ | ✓ |
| `content.posts.edit` | Edit post | | ✓ | ✓ | ✓ |
| `content.posts.publish` | Publish / unpublish post | | | ✓ | ✓ |
| `content.posts.delete` | Hapus post (soft delete) | | | | ✓ |
| `content.pages.view` | Lihat page | ✓ | ✓ | ✓ | ✓ |
| `content.pages.create` | Buat page | | ✓ | ✓ | ✓ |
| `content.pages.edit` | Edit page | | ✓ | ✓ | ✓ |
| `content.pages.publish` | Publish page | | | ✓ | ✓ |
| `content.pages.delete` | Hapus page | | | | ✓ |

### 4.3 Media

| Key | Label | Read only | Edit | Co-Admin |
|-----|-------|:---------:|:----:|:--------:|
| `media.view` | Lihat perpustakaan media | ✓ | ✓ | ✓ |
| `media.upload` | Upload file | | ✓ | ✓ |
| `media.delete` | Hapus media | | | ✓ |

### 4.4 SEO & Meta

| Key | Label | Read only | Edit | Co-Admin |
|-----|-------|:---------:|:----:|:--------:|
| `seo.view` | Lihat meta SEO | ✓ | ✓ | ✓ |
| `seo.edit` | Edit meta per konten / bulk SEO | | ✓ | ✓ |
| `seo.sitemap` | Generate / purge sitemap | | | ✓ |
| `seo.redirect` | Kelola redirect | | | ✓ |

Detail hierarki meta: [14-setup-meta-dan-seo.md](./14-setup-meta-dan-seo.md).

### 4.5 Operasi & sistem domain

| Key | Label | Read only | Edit | Co-Admin |
|-----|-------|:---------:|:----:|:--------:|
| `jobs.view` | Lihat status job domain | ✓ | ✓ | ✓ |
| `jobs.create` | Jalankan bulk / batch job | | | ✓ |
| `jobs.cancel` | Batalkan job milik sendiri | | | ✓ |
| `reports.view` | Lihat laporan domain | ✓ | ✓ | ✓ |

### 4.6 Larangan (tidak pernah via share)

| Key | Siapa saja yang boleh |
|-----|---------------------|
| `domain.transfer` | **Super Admin** saja |
| `domain.delete` | Owner + Super Admin |
| `host.manage` | **Super Admin** saja (subdomain produk) |
| `system.settings` | **Super Admin** saja |

---

## 5. Preset → JSON (contoh)

```json
{
  "preset": "read_only",
  "permissions": {
    "domain.view": true,
    "domain.settings.view": true,
    "content.posts.view": true,
    "content.pages.view": true,
    "media.view": true,
    "seo.view": true,
    "jobs.view": true,
    "reports.view": true
  }
}
```

```json
{
  "preset": "co_admin",
  "permissions": {
    "domain.view": true,
    "domain.settings.view": true,
    "domain.settings.edit": true,
    "domain.share": true,
    "content.posts.view": true,
    "content.posts.create": true,
    "content.posts.edit": true,
    "content.posts.publish": true,
    "content.posts.delete": true,
    "content.pages.view": true,
    "content.pages.create": true,
    "content.pages.edit": true,
    "content.pages.publish": true,
    "content.pages.delete": true,
    "media.view": true,
    "media.upload": true,
    "media.delete": true,
    "seo.view": true,
    "seo.edit": true,
    "seo.sitemap": true,
    "seo.redirect": true,
    "jobs.view": true,
    "jobs.create": true,
    "jobs.cancel": true,
    "reports.view": true
  }
}
```

---

## 6. Pengecekan di Backend (Go)

```go
func RequirePermission(userID, domainID int64, perm string) error {
  if user.IsSuperAdmin() { return nil }
  if access.IsOwner(userID, domainID) { return nil } // owner = semua kecuali system.*
  p := access.GetSharePermissions(userID, domainID)
  if !p[perm] { return ErrForbidden }
  return nil
}
```

| Skenario | Dampak jika tidak dicek |
|----------|------------------------|
| User read_only memanggil `POST /posts` | Kebocoran write — **wajib** middleware |
| Share tanpa `domain.share` mengundang user | Bypass approval chain |

Cache permission per request (in-memory map) — invalidasi saat share di-update.

---

## 7. Matriks RBAC Lengkap (Menu Admin)

| Menu / Aksi | Super Admin | Owner | Share (read_only) | Share (edit) | Share (co_admin) |
|-------------|:-----------:|:-----:|:-----------------:|:------------:|:----------------:|
| Dashboard global | ✓ | — | — | — | — |
| Dashboard domain | ✓ | ✓ | view | view | ✓ |
| Semua domain | ✓ | — | — | — | — |
| Domain milik | ✓ | ✓ | — | — | — |
| Domain di-share | ✓ | — | ✓ | ✓ | ✓ |
| Setup → Host | ✓ | — | — | — | — |
| Setup → Backend | ✓ | — | — | — | — |
| Setup → Meta (global) | ✓ | — | — | — | — |
| Berbagi akses + checklist | ✓ | ✓ | — | — | ✓ |
| Approve undangan co-admin | ✓ | ✓ | — | — | — |
| Transfer ownership | ✓ | — | — | — | — |
| Konten CRUD | ✓ | ✓ | view | create/edit | ✓ |
| Publish | ✓ | ✓ | — | — | ✓ |
| Media upload | ✓ | ✓ | — | ✓ | ✓ |
| SEO edit | ✓ | ✓ | view | ✓ | ✓ |
| Bulk / Jobs | ✓ | ✓ | — | — | ✓ |
| Kelola user sistem | ✓ | — | — | — | — |

---

## 8. Perubahan Database

```sql
ALTER TABLE domain_shares
  ADD COLUMN permission_preset TEXT,
  ADD COLUMN permissions JSONB NOT NULL DEFAULT '{}';

-- domain_share_invitations: sama, simpan preset + permissions untuk preview owner
ALTER TABLE domain_share_invitations
  ADD COLUMN permission_preset TEXT,
  ADD COLUMN permissions JSONB NOT NULL DEFAULT '{}';
```

`user_domain_access` tetap untuk list cepat; optional column `permission_summary` (preset name) untuk UI badge.

---

## 9. Dokumen Terkait

- Share + approval → [09](./09-model-domain-host-dan-subdomain.md)
- Login & session → [12](./12-autentikasi-dan-login-aman.md)
- Setup backend → [13](./13-setup-backend-dan-sistem.md)
- Meta → [14](./14-setup-meta-dan-seo.md)
