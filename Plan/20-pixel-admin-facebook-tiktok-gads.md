# 20 — Pixel Hub: Kolaborasi dengan Facebook, TikTok & Google Ads

> Modul **`/admin/pixel/*`** bukan sekadar form “tempel Pixel ID”. Ini adalah **wadah kolaborasi** antara tim Anda dan platform iklan — setara peran **Stape.io**, **GTM Server-Side**, atau **Meta CAPI Gateway**, tetapi **native** di CMS Seosementara, dioptimalkan untuk **ribuan domain** dari satu panel.  
> Selaras: [13](./13-setup-backend-dan-sistem.md), [11](./11-rbac-dan-permission-share.md), [14](./14-setup-meta-dan-seo.md), [15](./15-setup-cloudflare-integrasi.md), [19](./19-modul-url-shortlink.md).

---

## 1. Apa yang Dimaksud “Pixel” (Pemahaman Bersama)

**Pixel** (Meta / TikTok) atau **tag konversi** (Google Ads) adalah sistem **pelacakan event** — kunjungan, klik, lead, pembelian — agar algoritma iklan bisa **optimasi** (CPA lebih murah, audience lebih tepat).

| Istilah | Arti praktis |
|---------|----------------|
| **Browser pixel** | Skrip JS (`fbq`, `ttq`, `gtag`) di browser pengunjung → kirim ke platform |
| **Server-side / CAPI** | Server Anda kirim event ke API platform (Meta CAPI, TikTok Events API, Google Enhanced Conversions) |
| **Event** | Satu aksi terukur: `PageView`, `Lead`, `Purchase`, … |
| **Dedup** | Meta/Google menghindari hitung ganda jika browser + server kirim event sama (`event_id`, `fbc`/`fbp`) |
| **First-party** | Request tracking dari **domain Anda** (`pelacak.seosementara.org`), bukan `connect.facebook.net` — lebih sulit diblokir adblock |

**Mengapa orang memakai layanan pihak ketiga (Stape, GTM SS, CAPIG)?**  
Bukan karena “pixel pintar” secara ajaib — melainkan karena mereka menyelesaikan **enam masalah operasional** di bawah. **Pixel Hub CMS** hadir untuk menggantikan kebutuhan berlangganan itu **tanpa** kehilangan kontrol data.

---

## 2. Visi: Pixel Hub sebagai Wadah Kolaborasi

```mermaid
flowchart TB
  subgraph visitor [Pengunjung]
    Site[Situs / shortlink / landing]
  end
  subgraph firstparty [First-party - domain Anda]
    JS[sseo-track.js ~3KB]
    Collect[POST /collect]
  end
  subgraph hub [Pixel Hub - Mini CPU]
    Ingest[Ingest + validasi + consent]
    Queue[Antrian pixel_dispatch]
    Map[Event catalog - 1 definisi N platform]
    Privacy[Hash PII + filter bot]
    Dispatch[Fan-out CAPI]
  end
  subgraph admin [Admin - Kolaborasi Tim]
    Overview[/admin/pixel]
    FB[/admin/pixel/facebook]
    TT[/admin/pixel/tiktok]
    GA[/admin/pixel/gads]
    Catalog[/admin/pixel/events]
  end
  subgraph platforms [Platform Iklan]
    Meta[Meta CAPI + Events Manager]
    TikTok[TikTok Events API]
    Google[Google Ads + GA4 MP]
  end
  Site --> JS
  JS --> Collect
  Collect --> Ingest
  Ingest --> Queue
  Queue --> Map
  Map --> Privacy
  Privacy --> Dispatch
  Dispatch --> Meta
  Dispatch --> TikTok
  Dispatch --> Google
  Overview --> hub
  FB --> Meta
  TT --> TikTok
  GA --> Google
```

| Peran halaman | Bukan hanya… | Melainkan… |
|---------------|--------------|------------|
| `/admin/pixel/` | Dashboard kosong | **Status hub**: recovery rate vs adblock, antrian, error CAPI, perbandingan kanal |
| `/admin/pixel/facebook/` | Form Pixel ID | **Ruang kerja Meta**: koneksi CAPI, test Events Manager, diagnosa match rate, assign ribuan domain |
| `/admin/pixel/tiktok/` | Idem | **Ruang kerja TikTok** |
| `/admin/pixel/gads/` | Idem | **Ruang kerja Google** |
| `/admin/pixel/events/` | - | **Katalog event** — definisi sekali, terjemah ke FB + TT + GAds |

---

## 3. Enam Masalah Industri → Solusi Native di CMS

| # | Masalah (yang Anda jelaskan) | Solusi Pixel Hub |
|---|------------------------------|------------------|
| **1** | Ad-blocker & iOS (ITP) memblokir `connect.facebook.net` | **First-party endpoint**: `pelacak.seosementara.org` (atau `t.{domain-klien}`) via Cloudflare DNS + Tunnel/Worker [15]. Browser hanya memuat **`sseo-track.js`** dari domain sendiri |
| **2** | CAPI rumit — butuh backend, queue, maintenance | **Built-in**: Go service + worker `pixel_dispatch` + retry + dead-letter. Admin: **aktifkan CAPI** tanpa coding — sama seperti “klik di Stape” |
| **3** | Ribuan domain — edit pixel satu-satu mustahil | **Mass assign** dari `/admin/pixel/*/domains` + policy template (grup domain → satu set pixel). Ubah config sekali → propagasi ke semua domain terikat |
| **4** | Banyak skrip pixel = PageSpeed anjlok | Mode default **`server_first`**: satu skrip ringan di browser; **fan-out** ke FB/TT/Google di server. Mode `hybrid` opsional untuk dedup browser+CAPI |
| **5** | Privasi GDPR — data mentah ke Meta | **Privacy gateway** di server: hash SHA256 email/telp, strip field sensitif, consent flag, log audit sebelum kirim |
| **6** | Empat platform = empat logika event | **Event catalog**: definisi `purchase` sekali → mapping otomatis ke `Purchase` (FB), `CompletePayment` (TT), `purchase` (Google) |

### Positioning vs layanan pihak ketiga

| Layanan eksternal | Yang mereka jual | Yang CMS lakukan sendiri |
|-------------------|------------------|---------------------------|
| Stape.io / GTM Server-Side | First-party + CAPI hosting | Subdomain `pelacak.*` + Tunnel + hub Go |
| Meta CAPIG | CAPI tanpa kode | Tab Setup Facebook + token di `pixel_credentials` |
| Segment / CDP | Satu event → banyak destinasi | Tabel `pixel_event_definitions` + `pixel_platform_mappings` |
| Cloudflare Zaraz | Tag di edge | **Opsional** fase 3 — bisa dipakai bersamaan atau diganti hub Go |

**Keuntungan native:** tidak ada biaya per-event pihak ketiga, data tetap di PostgreSQL Anda, mass deploy selaras modul domain portfolio.

---

## 4. Mode Operasi (Pilih di `/admin/pixel/hub/settings`)

| Mode | Browser | Server | Kapan dipakai |
|------|---------|--------|---------------|
| **`server_first`** *(disarankan)* | Hanya `sseo-track.js` → POST collect | CAPI ke semua platform aktif | Skala besar, SEO, adblock tinggi |
| **`hybrid`** | `sseo-track.js` + optional `fbq`/`ttq`/`gtag` tipis | CAPI + dedup `event_id` | Perlu sinyal browser untuk EMQ Meta |
| **`legacy_client`** | Skrip penuh platform di `<head>` | Opsional backup CAPI | Migrasi dari setup lama |

Toggle per **scope** (global / domain portfolio / shortlink).

---

## 5. Arsitektur First-Party & Collect API

### 5.1 Hostname pelacakan

| Pola | Contoh | Catatan |
|------|--------|---------|
| Subdomain produk | `pelacak.seosementara.org` | Satu untuk semua domain portfolio |
| Per domain klien (fase 2) | `t.rezekibelanja.com` | CNAME → Cloudflare → Worker/Tunnel |
| Path di apex | `seosementara.org/t/collect` | Fallback jika subdomain belum siap |

Setup DNS & route dari [15-setup-cloudflare](./15-setup-cloudflare-integrasi.md) tab **Pixel / Tracking**.

### 5.2 Alur satu event

1. Pengunjung membuka halaman → `sseo-track.js` kirim `POST /collect` dengan `event: page_view`, `url`, `domain_id`, `session_id`, `_fbp`/`_fbc` jika ada cookie.
2. **Ingest** (Go): validasi origin, bot score CF (header), cek consent.
3. Tulis `pixel_events` status `pending`.
4. Worker **`pixel_dispatch`**: enrich user_data (IP, UA), hash PII, map ke platform, kirim batch API.
5. Update status `sent` / `failed` + `platform_event_id`.

**Shortlink [19]:** event `click` hanya enqueue — **tidak** blocking redirect HTTP 302.

### 5.3 Skrip browser (ringan)

```html
<script async src="https://pelacak.seosementara.org/sseo-track.js"
        data-site="{{.SiteKey}}"
        data-consent="{{.ConsentRequired}}"></script>
```

```javascript
// sseo-track.js - konsep
window.sseo = window.sseo || { q: [] };
window.sseo.track = (name, props) => {
  navigator.sendBeacon('/collect', JSON.stringify({ event: name, ...props, ts: Date.now() }));
};
sseo.track('page_view', { path: location.pathname });
```

**Tidak** memuat `connect.facebook.net` di mode `server_first`.

---

## 6. Event Catalog — Satu Sumber, Banyak Platform

Halaman **`/admin/pixel/events/`** (shared, bukan duplikat di tiap kanal).

| Kolom katalog | Contoh |
|---------------|--------|
| `canonical_name` | `purchase` |
| `label_id` | `evt_purchase_01` |
| Trigger | `manual`, `shortlink_click`, `form_submit`, `webhook` |
| Aktif di platform | checkbox FB / TT / GAds |

Tabel mapping:

```sql
CREATE TABLE pixel_event_definitions (
  id              BIGSERIAL PRIMARY KEY,
  canonical_name  TEXT NOT NULL UNIQUE,
  label           TEXT NOT NULL,
  description     TEXT,
  trigger_type    TEXT NOT NULL,
  is_active       BOOLEAN NOT NULL DEFAULT true,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pixel_platform_mappings (
  id                    BIGSERIAL PRIMARY KEY,
  event_definition_id   BIGINT NOT NULL REFERENCES pixel_event_definitions(id),
  platform              TEXT NOT NULL CHECK (platform IN ('facebook','tiktok','gads')),
  platform_event_name   TEXT NOT NULL,
  extra_params          JSONB NOT NULL DEFAULT '{}',
  UNIQUE (event_definition_id, platform)
);
```

Contoh baris mapping untuk `purchase`:

| Platform | `platform_event_name` |
|----------|----------------------|
| facebook | `Purchase` |
| tiktok | `CompletePayment` |
| gads | `purchase` (+ conversion label di `extra_params`) |

Halaman **facebook/tiktok/gads** menampilkan mapping **read-only** + link “Edit di Event Catalog” — kolaborasi terpusat, konfigurasi per-platform hanya untuk **credential & pixel ID**.

---

## 7. Struktur Menu Admin (Diperbarui)

```
/admin/pixel/
├── overview              → KPI hub: events/hari, gagal CAPI, recovery vs blocked estimate
├── hub/
│   ├── settings          → Mode server_first / hybrid, hostname pelacak, consent
│   ├── privacy           → Hash rules, field allowlist, retention
│   └── deploy            → Mass deploy snippet ke domain (batch job)
├── events/               → Event catalog + mapping 3 platform
│
├── facebook/             → Kolaborasi Meta
│   ├── setup             → Pixel ID, CAPI token, test event code
│   ├── connection        → Status token, EMQ score, link Events Manager
│   ├── domains           → Assign + template mass
│   ├── diagnostics       → Failed events, dedup, adblock recovery %
│   └── analytics         → Internal + sync API
│
├── tiktok/               → Pola sama
└── gads/                 → Pola sama (+ tab GA4)
```

Sidebar: **Pixel** → Overview | Event Catalog | Facebook | TikTok | Google Ads.

---

## 8. Halaman Kolaborasi per Platform

Setiap `/admin/pixel/{facebook|tiktok|gads}/` adalah **ruang kerja** dengan platform tersebut — bukan sekadar CRUD.

### 8.1 Tab Setup & Connection

| Fitur | Manfaat kolaborasi |
|-------|-------------------|
| Pixel / Conversion ID | Sinkron dengan Events Manager / TikTok Ads / Google Ads UI |
| Credential CAPI / OAuth | Disimpan terenkripsi; tombol **Test koneksi** |
| **Test Event** | Kirim event uji → tampil di Events Manager (kode test) |
| **Buka di platform** | Deep link ke dashboard eksternal untuk verifikasi domain |
| Status badge | `connected` / `token_expired` / `rate_limited` |

### 8.2 Tab Domains (mass)

| Aksi | Perilaku |
|------|----------|
| Assign pixel ke 1 domain | Dropdown domain portfolio |
| Assign ke grup | Filter tag owner / grup bisnis |
| **Propagasi** | Job `pixel_deploy_snippet` update meta domain tanpa edit kode manual |
| Preview | “Simulasi event” untuk domain terpilih |

### 8.3 Tab Diagnostics

| Metrik | Interpretasi |
|--------|--------------|
| Events `pending` > 5 menit | Worker backlog — scale batch |
| `failed` rate | Token salah / payload invalid |
| Browser vs server count | Estimasi data yang “hilang” sebelum hub |
| Recovery rate | Setelah first-party — target naik vs baseline legacy |

### 8.4 Tab Analytics

Kombinasi **log internal** (`pixel_events`) + **sync berkala** API platform (terbatas di tier gratis — tampilkan disclaimer + link eksternal).

---

## 9. Privacy Gateway (Satpam Data)

Sebelum fan-out ke Meta/TikTok/Google:

| Langkah | Aturan |
|---------|--------|
| Consent | Jika `consent_required`: hanya kirim `marketing` setelah cookie/HTMX banner |
| Normalisasi email | lowercase, trim → SHA256 |
| Telepon | E.164 → SHA256 |
| Strip | Jangan kirim nama lengkap mentah kecuali Enhanced Conversions Google (hash wajib) |
| Bot | CF `cf-bot-score` > threshold → mark `dropped_bot` |
| Audit | `pixel_privacy_log` (opsional) — siapa kirim apa, tanpa PII mentah |

Halaman: `/admin/pixel/hub/privacy`.

---

## 10. Schema Database (Lengkap)

```sql
CREATE TABLE pixel_hub_settings (
  id                    BIGSERIAL PRIMARY KEY,
  tracking_hostname     TEXT NOT NULL DEFAULT 'pelacak.seosementara.org',
  default_mode          TEXT NOT NULL DEFAULT 'server_first'
                        CHECK (default_mode IN ('server_first','hybrid','legacy_client')),
  consent_required      BOOLEAN NOT NULL DEFAULT false,
  collect_path          TEXT NOT NULL DEFAULT '/collect',
  script_version        TEXT NOT NULL DEFAULT '1',
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pixel_configs (
  id                BIGSERIAL PRIMARY KEY,
  platform          TEXT NOT NULL CHECK (platform IN ('facebook','tiktok','gads')),
  scope             TEXT NOT NULL CHECK (scope IN ('global','managed_domain','shortlink')),
  managed_domain_id BIGINT REFERENCES managed_domains(id) ON DELETE CASCADE,
  name              TEXT NOT NULL,
  is_active         BOOLEAN NOT NULL DEFAULT true,
  mode_override     TEXT CHECK (mode_override IN ('server_first','hybrid','legacy_client')),
  external_ids      JSONB NOT NULL DEFAULT '{}',
  credentials_ref   BIGINT REFERENCES pixel_credentials(id),
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pixel_credentials (
  id                BIGSERIAL PRIMARY KEY,
  platform          TEXT NOT NULL,
  name              TEXT NOT NULL,
  secret_ciphertext BYTEA NOT NULL,
  secret_nonce      BYTEA NOT NULL,
  last_validated_at TIMESTAMPTZ,
  validation_status TEXT,
  updated_by        BIGINT REFERENCES users(id),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE pixel_events (
  id                BIGSERIAL PRIMARY KEY,
  canonical_event   TEXT,
  platform          TEXT,  -- NULL = belum di-map; terisi per baris fan-out
  pixel_config_id   BIGINT REFERENCES pixel_configs(id),
  event_name        TEXT NOT NULL,
  event_id          TEXT,   -- dedup UUID
  managed_domain_id BIGINT,
  url_link_id       BIGINT,
  payload           JSONB NOT NULL,
  status            TEXT NOT NULL DEFAULT 'pending',
  platform_event_id TEXT,
  error_message     TEXT,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_pixel_events_status_created
  ON pixel_events (status, created_at)
  WHERE status = 'pending';
```

---

## 11. Worker & Performa (Skala Ribuan Domain)

| Job | Fungsi | Batasan |
|-----|--------|---------|
| `pixel_dispatch` | Kirim batch pending → platform APIs | Batch 50–100, interval 5–10 dtk |
| `pixel_deploy_snippet` | Update flag/meta domain untuk inject script | Chunk 200 domain/job |
| `pixel_sync_analytics` | Pull ringkasan API (harian) | Rate limit API platform |

**Prinsip skala:** tidak ada loop 3000 domain di request HTTP — semua propagasi via **job queue** [13].

---

## 12. RBAC [11]

| Permission | Label |
|------------|-------|
| `pixel.view` | Overview + diagnostics read |
| `pixel.hub.manage` | Settings, privacy, deploy mass |
| `pixel.events.manage` | Event catalog |
| `pixel.facebook.manage` | Credential + test Meta |
| `pixel.tiktok.manage` | Credential + test TikTok |
| `pixel.gads.manage` | Credential + test Google |
| `pixel.analytics` | Semua tab analytics |

---

## 13. API (Ringkas)

| Method | Path | Fungsi |
|--------|------|--------|
| POST | `/collect` | Ingest first-party (publik, rate limit) |
| GET | `/sseo-track.js` | Skrip ringan (cache CF) |
| GET/PATCH | `/api/admin/pixel/hub/settings` | Mode & hostname |
| GET/POST | `/api/admin/pixel/events` | Katalog |
| POST | `/api/admin/pixel/hub/deploy` | Mass deploy job |
| GET/PATCH | `/api/admin/pixel/facebook/setup` | Meta credential |
| POST | `/api/admin/pixel/facebook/test` | Test CAPI |
| GET | `/api/admin/pixel/facebook/diagnostics` | Errors & recovery |
| GET | `/api/admin/pixel/facebook/analytics` | Agregat |
| * | `/api/admin/pixel/tiktok/*` | Paralel TikTok |
| * | `/api/admin/pixel/gads/*` | Paralel Google |

---

## 14. Integrasi Cloudflare [15]

| Komponen | Peran |
|----------|-------|
| DNS `pelacak.*` | First-party hostname |
| Tunnel / Worker route | `/collect` → Go ingest |
| Cache `sseo-track.js` | Edge cache, TTL panjang + purge saat versi naik |
| Bot Management | Filter sebelum enqueue |
| Zaraz (opsional fase 4) | Alternatif edge — hub Go tetap sumber kebenaran event internal |

---

## 15. Mapping Aksi CMS → Event

| Aksi CMS | Canonical | FB | TikTok | Google |
|----------|-----------|-----|--------|--------|
| Shortlink klik [19] | `click` | `ViewContent` | `Click` | conversion click |
| Page view | `page_view` | `PageView` | `Pageview` | `page_view` |
| Form lead | `lead` | `Lead` | `SubmitForm` | `generate_lead` |
| Purchase | `purchase` | `Purchase` | `CompletePayment` | `purchase` |

Semua bisa di-override di Event Catalog tanpa ubah kode Go.

---

## 16. Roadmap Implementasi

| Fase | Deliverable | Mengatasi masalah # |
|------|-------------|---------------------|
| **MVP** | `sseo-track.js` + `/collect` + hub settings + Facebook CAPI + `pixel_dispatch` | 1, 2 |
| **Fase 2** | Event catalog + mapping TikTok & GAds + halaman diagnostics | 6 |
| **Fase 3** | Mass deploy domains + privacy hash + consent banner | 3, 5 |
| **Fase 4** | Analytics sync + hybrid dedup + per-domain `t.{domain}` | 4, 1 |
| **Fase 5** | Cloudflare Zaraz opsi / EMQ tuning | 1, 4 |

---

## 17. Checklist

- [ ] `pixel_hub_settings` + hostname `pelacak.*` di CF
- [ ] Endpoint `POST /collect` + `sseo-track.js`
- [ ] Worker `pixel_dispatch` + retry
- [ ] UI `/admin/pixel/` overview + `/admin/pixel/hub/settings`
- [ ] UI kolaborasi `/admin/pixel/facebook/*` (setup, connection, diagnostics)
- [ ] Event catalog `/admin/pixel/events/`
- [ ] Mass deploy job
- [ ] Privacy gateway (hash + consent)
- [ ] TikTok & GAds fan-out
- [ ] RBAC permissions

---

## 18. Dokumen Terkait

- [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- [08-roadmap-implementasi.md](./08-roadmap-implementasi.md)
- [11-rbac-dan-permission-share.md](./11-rbac-dan-permission-share.md)
- [15-setup-cloudflare-integrasi.md](./15-setup-cloudflare-integrasi.md)
- [19-modul-url-shortlink.md](./19-modul-url-shortlink.md)
