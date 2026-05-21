# 20 — Pixel Admin (Facebook, TikTok, Google Ads)

> Halaman admin **`/admin/pixel/*`** untuk **setup**, **komunikasi** ke platform iklan, dan **analitik** per kanal.  
> Selaras: [13-setup-backend](./13-setup-backend-dan-sistem.md), [11-rbac](./11-rbac-dan-permission-share.md), [14-setup-meta](./14-setup-meta-dan-seo.md), [19-modul-url-shortlink](./19-modul-url-shortlink.md).

## 1. Tujuan

| Tujuan | Keterangan |
|--------|------------|
| Satu tempat kelola pixel | Admin tidak bolak-balik ke Business Manager / Google / TikTok untuk rutinitas |
| Tiga kanal terpisah | UI & credential terpisah — tidak dicampur |
| **Komunikasi dua arah** | Client snippet + **server-side API** (Events API / CAPI / Google) |
| Skala domain | Ribuan domain portfolio — pixel global, per domain, atau per grup |
| Analisis | Dashboard per kanal di halaman masing-masing |

---

## 2. Struktur Menu & URL Admin

```
/admin/pixel/                    → ringkasan semua kanal (kartu FB / TT / GAds)
│
├── /admin/pixel/facebook/       → Modul Meta (Facebook) Pixel
│   ├── setup                    → Pixel ID, CAPI token, domain verify
│   ├── events                   → log event terkirim + test event
│   ├── domains                  → assign pixel ke domain portfolio
│   └── analytics                → grafik & metrik (sync API + internal)
│
├── /admin/pixel/tiktok/         → Modul TikTok Pixel
│   ├── setup
│   ├── events
│   ├── domains
│   └── analytics
│
└── /admin/pixel/gads/           → Modul Google Ads / Google tag
    ├── setup                    → Conversion ID, labels, GA4 link
    ├── events
    ├── domains
    └── analytics
```

Sidebar admin: grup **Pixel** dengan tiga submenu (ikon + badge status koneksi hijau/merah).

```mermaid
flowchart TB
  subgraph admin [Admin Panel]
    P[/admin/pixel]
    FB[/admin/pixel/facebook]
    TT[/admin/pixel/tiktok]
    GA[/admin/pixel/gads]
  end
  subgraph platforms [Platform APIs]
    Meta[Meta Graph / CAPI]
    TikTok[TikTok Events API]
    Google[Google Ads / GA4 Data API]
  end
  subgraph emit [Sumber Event]
    Web[Snippet di halaman publik]
    URL[Redirect shortlink url.*]
    Srv[Backend server-side]
  end
  P --> FB
  P --> TT
  P --> GA
  FB --> Meta
  TT --> TikTok
  GA --> Google
  Web --> Srv
  URL --> Srv
  Srv --> Meta
  Srv --> TikTok
  Srv --> Google
```

---

## 3. Konsep: Cara Pixel “Berkomunikasi”

| Lapisan | Fungsi | Platform |
|---------|--------|----------|
| **Browser (client)** | Script pixel di `<head>` — PageView, klik | FB `fbq`, TikTok `ttq`, Google `gtag` |
| **Server (server-side)** | Event API — lebih akurat, tidak terblokir adblock | Meta CAPI, TikTok Events API, Google Enhanced Conversions |
| **Internal CMS** | Log event + korelasi domain / shortlink | Tabel `pixel_events` |

**Prinsip:** Admin panel menyimpan **credential + pixel ID** → backend **mengirim** event ke platform + menampilkan **status & statistik** dari sync API.

---

## 4. Scope Penempatan Pixel

| Scope | Contoh | Dipakai saat |
|-------|--------|--------------|
| **Global produk** | Satu pixel Seosementara | `seosementara.org`, subdomain produk |
| **Per domain portfolio** | Pixel klien `rezekibelanja.com` | Preview / landing / (fase 2) situs publik domain |
| **Per shortlink** | Event `Click` saat redirect [19] | `url.seosementara.org/...` |
| **Per kampanye ads** | Modul `ads.*` [18] | Landing kampanye |

```sql
CREATE TABLE pixel_configs (
  id              BIGSERIAL PRIMARY KEY,
  platform        TEXT NOT NULL CHECK (platform IN ('facebook','tiktok','gads')),
  scope           TEXT NOT NULL CHECK (scope IN ('global','managed_domain','shortlink')),
  managed_domain_id BIGINT REFERENCES managed_domains(id) ON DELETE CASCADE,
  name            TEXT NOT NULL,
  is_active       BOOLEAN NOT NULL DEFAULT true,
  -- platform-specific IDs (non-secret)
  external_ids    JSONB NOT NULL DEFAULT '{}',
  -- secrets encrypted
  credentials_ref TEXT,  -- pointer ke pixel_credentials
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Contoh `external_ids`:

```json
// facebook
{ "pixel_id": "1234567890", "business_id": "..." }

// tiktok
{ "pixel_code": "CXXXX", "advertiser_id": "..." }

// gads
{ "conversion_id": "AW-xxx", "conversion_label": "abc123", "ga4_measurement_id": "G-XXXX" }
```

---

## 5. Kredensial (Terenkripsi)

```sql
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
```

| Platform | Secret yang disimpan | Validasi di Setup |
|----------|----------------------|-------------------|
| **Facebook** | CAPI **Access Token** (+ optional app secret) | Test event ke Meta |
| **TikTok** | Events API **Access Token** | Test event |
| **Google Ads** | OAuth refresh token / service account JSON | List conversions |

Sama seperti Cloudflare [15]: **tidak** tampilkan secret penuh di UI — masked + tombol **Test koneksi**.

---

## 6. Halaman `/admin/pixel/facebook/`

### 6.1 Setup

| Field | Keterangan |
|-------|------------|
| Pixel ID | Dari Events Manager Meta |
| CAPI Access Token | Server-side |
| Test Event Code | Mode debug Meta |
| Aktifkan CAPI | Toggle |
| Aktifkan browser pixel | Toggle snippet |
| Domain portfolio | Multi-select / search assign |

Tombol:

- **Simpan & uji koneksi** → `POST .../facebook/test`
- **Kirim event uji** → `PageView` test ke Meta

### 6.2 Events (log internal)

| Kolom | Sumber |
|-------|--------|
| Waktu | `pixel_events` |
| Event name | `PageView`, `Lead`, `Purchase`, … |
| Domain / link | `managed_domain_id` / `url_link_id` |
| Status kirim | `sent`, `failed`, `pending` |
| Platform response | `event_id` Meta / error |

### 6.3 Analytics

| Metrik | Sumber |
|--------|--------|
| Event terkirim 7/30 hari | Internal DB |
| Match rate CAPI vs browser (jika ada) | Meta API |
| Top domain by events | Agregat CMS |
| Sync dari Meta | Job harian (Insights API terbatas — gunakan Events Manager export atau Marketing API) |

**Catatan:** Meta tidak memberikan semua data real-time di API gratis — kombinasikan **internal log** + **sync berkala** + link ke Events Manager (buka eksternal).

### 6.4 Komunikasi teknis (Facebook)

| Arah | API |
|------|-----|
| CMS → Meta | [Conversions API](https://developers.facebook.com/docs/marketing-api/conversions-api) `POST /{pixel-id}/events` |
| Meta → CMS | Webhook (opsional fase 2) — jarang dipakai |
| Browser → Meta | Snippet `fbq` diinjeksi dari template |

Payload minimal server-side:

```json
{
  "event_name": "PageView",
  "event_time": 1716300000,
  "action_source": "website",
  "event_source_url": "https://rezekibelanja.com/",
  "user_data": { "client_ip_address": "...", "client_user_agent": "..." }
}
```

---

## 7. Halaman `/admin/pixel/tiktok/`

### 7.1 Setup

| Field | Keterangan |
|-------|------------|
| Pixel Code | TikTok Events Manager |
| Events API Access Token | Server-side |
| Aktifkan browser pixel | `ttq` snippet |

### 7.2 Analytics & Events

Struktur sama dengan Facebook — tabel `pixel_events` dengan `platform = tiktok`.

### 7.3 Komunikasi teknis (TikTok)

| Arah | API |
|------|-----|
| CMS → TikTok | Events API 2.0 — `POST` event batch |
| Browser → TikTok | TikTok Pixel JS |

Event contoh: `ViewContent`, `ClickButton`, `CompletePayment` — mapping dari aksi CMS (publish, shortlink click).

---

## 8. Halaman `/admin/pixel/gads/`

### 8.1 Setup

| Field | Keterangan |
|-------|------------|
| Google Ads Conversion ID | `AW-xxxxxxxx` |
| Conversion label(s) | Per tipe konversi |
| GA4 Measurement ID | `G-xxxxxxxx` (opsional, analitik) |
| Google Tag ID | Untuk gtag.js |
| OAuth / Service Account | Untuk upload konversi server-side |

### 8.2 Analytics

| Metrik | Sumber |
|--------|--------|
| Conversions (uploaded) | Internal + Google Ads API |
| Clicks vs conversions | Google Ads reporting sync |
| Per domain | Label suffix atau custom parameters |

### 8.3 Komunikasi teknis (Google)

| Arah | API / metode |
|------|----------------|
| Browser | gtag.js — `gtag('config', 'AW-...')`, `gtag('event', 'conversion', {...})` |
| Server | **Enhanced Conversions** / Offline conversion upload / GA4 Measurement Protocol |
| CMS → Google | Google Ads API `uploadClickConversions` atau GA4 `mp/collect` |

**`gads`** di menu = **Google Ads conversions** + **Google tag (gtag)** — bisa satu halaman dengan tab **Ads** | **GA4**.

---

## 9. Injeksi Snippet (Browser Pixel)

### 9.1 Dari mana script dimuat

| Target | Cara |
|--------|------|
| Apex / subdomain produk | Partial `<head>` dari [14] — `meta.global.pixels` |
| Domain portfolio (preview/publik) | `managed_domain_meta.pixels` |
| Shortlink redirect page | Opsional intermediate page 200ms + fire event |

Backend render:

```html
{{if .Pixels.Facebook.Active}}
<script>!function(f,b,e,v,n,t,s){...}(...);
fbq('init', '{{.Pixels.Facebook.PixelID}}');
fbq('track', 'PageView');
</script>
{{end}}
```

TikTok & gtag — template terpisah per platform di `Frontend-admin` tidak dipakai untuk publik; publik di `Frontend-Users` atau HTML dari Go.

### 9.2 Consent (disarankan fase 2)

| Mode | Perilaku |
|------|----------|
| Tanpa consent | Fire langsung (MVP internal) |
| Dengan consent | HTMX banner → load pixel setelah setuju |

---

## 10. Event Internal → Platform (Mapping)

| Aksi di CMS | Facebook | TikTok | Google Ads |
|-------------|----------|--------|------------|
| Shortlink klik [19] | `ViewContent` / custom | `Click` | `conversion` click |
| Publish post | `PageView` (jika URL live) | `ViewContent` | - |
| Form lead | `Lead` | `SubmitForm` | `generate_lead` |
| Purchase (fase 2) | `Purchase` | `CompletePayment` | `purchase` |

```sql
CREATE TABLE pixel_events (
  id                BIGSERIAL PRIMARY KEY,
  platform          TEXT NOT NULL,
  pixel_config_id   BIGINT REFERENCES pixel_configs(id),
  event_name        TEXT NOT NULL,
  managed_domain_id BIGINT,
  url_link_id       BIGINT,
  payload           JSONB NOT NULL,
  status            TEXT NOT NULL DEFAULT 'pending',
  platform_event_id TEXT,
  error_message     TEXT,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_pixel_events_platform_time
  ON pixel_events (platform, created_at DESC);
```

Worker job: `pixel_dispatch` — batch kirim `pending` ke Meta/TikTok/Google.

---

## 11. RBAC [11]

| Permission | Label |
|------------|-------|
| `pixel.view` | Lihat semua halaman pixel |
| `pixel.facebook.manage` | Setup + kirim test Facebook |
| `pixel.tiktok.manage` | Setup TikTok |
| `pixel.gads.manage` | Setup Google |
| `pixel.analytics` | Lihat analitik semua kanal |

| Role | Akses |
|------|-------|
| Super Admin | Semua |
| Platform Manager | `pixel.*` tanpa credential delete global (opsional) |
| Pekerja | Hanya pixel domain milik + share (scope `managed_domain`) |

---

## 12. API Admin (Ringkas)

| Method | Path |
|--------|------|
| GET | `/api/admin/pixel/overview` |
| GET/PATCH | `/api/admin/pixel/facebook/setup` |
| POST | `/api/admin/pixel/facebook/test` |
| POST | `/api/admin/pixel/facebook/events/test` |
| GET | `/api/admin/pixel/facebook/analytics` |
| GET | `/api/admin/pixel/facebook/events` |
| GET/PATCH | `/api/admin/pixel/tiktok/setup` |
| POST | `/api/admin/pixel/tiktok/test` |
| GET | `/api/admin/pixel/tiktok/analytics` |
| GET/PATCH | `/api/admin/pixel/gads/setup` |
| POST | `/api/admin/pixel/gads/test` |
| GET | `/api/admin/pixel/gads/analytics` |
| GET/PATCH | `/api/admin/pixel/domains/{domain_id}` | Assign configs per domain |

---

## 13. UI HTMX [17]

| Pola | Pemakaian |
|------|-----------|
| Tab dalam platform | Setup | Events | Domains | Analytics |
| Test connection | `hx-post` → swap alert success/error |
| Grafik | Server return partial `<canvas>` data atau library ringan |
| Assign domain | Combobox search + checklist pixel aktif |
| Status koneksi | Badge di sidebar submenu |

---

## 14. Integrasi Cloudflare [15]

| Fitur CF | Manfaat pixel |
|----------|---------------|
| Zaraz (opsional) | Ganti inject manual — kelola FB/Google/TikTok di edge |
| Bot score | Filter event bot sebelum kirim CAPI |
| Analytics | Banding traffic CF vs event pixel |

Setup opsional di `/admin/setup/cloudflare/` tab **Zaraz** — fase 3. MVP: inject dari Go/Pages.

---

## 15. Skenario & Dampak

| # | Skenario | Dampak |
|---|----------|--------|
| P1 | Token CAPI expired | Event gagal — banner merah di Setup |
| P2 | 3000 domain, pixel berbeda | Banyak `pixel_configs` — index `managed_domain_id` |
| P3 | Kirim event sync blocking redirect | Shortlink lambat — **async worker** wajib |
| P4 | Pekerja lihat pixel domain lain | Kebocoran — filter owner/share |
| P5 | Duplikat PageView browser+CAPI | Meta dedup pakai `event_id` + `fbc`/`fbp` |
| P6 | GDPR | Tanpa consent — risiko hukum EU |

---

## 16. Roadmap

| Fase | Deliverable |
|------|-------------|
| MVP | UI 3 halaman Setup + credential + test event Facebook saja |
| Fase 2 | TikTok + GAds setup + `pixel_events` log |
| Fase 3 | Analytics sync + per-domain assign + shortlink events |
| Fase 4 | Consent banner + Cloudflare Zaraz opsi |

---

## 17. Checklist Implementasi

- [ ] Menu sidebar Pixel + 3 submenu
- [ ] Tabel `pixel_configs`, `pixel_credentials`, `pixel_events`
- [ ] Halaman `/admin/pixel/facebook/setup` + test CAPI
- [ ] Worker `pixel_dispatch`
- [ ] Snippet inject hook di template publik
- [ ] RBAC permissions
- [ ] Ulang untuk TikTok & GAds

---

## 18. Dokumen Terkait

- [03-menu-dan-modul-cms.md](./03-menu-dan-modul-cms.md)
- [11-rbac-dan-permission-share.md](./11-rbac-dan-permission-share.md)
- [13-setup-backend-dan-sistem.md](./13-setup-backend-dan-sistem.md)
- [14-setup-meta-dan-seo.md](./14-setup-meta-dan-seo.md)
- [15-setup-cloudflare-integrasi.md](./15-setup-cloudflare-integrasi.md)
- [17-kontrak-htmx-dan-komponen-ui.md](./17-kontrak-htmx-dan-komponen-ui.md)
- [19-modul-url-shortlink.md](./19-modul-url-shortlink.md)
