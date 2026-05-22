# 25 — Data Pixel Selaras Meta (2026) — Bukan Hanya IP & Device

> **Diverifikasi ulang** terhadap dokumentasi resmi Meta (Conversions API / Parameters), diperbarui **Mei 2026**.  
> Dokumen ini menggantikan daftar data internal yang kurang lengkap.  
> CAPI: [23](./23-meta-conversions-api-kedalaman.md) · BM: [24](./24-meta-akun-bm-pixel-dan-optimasi-iklan.md)

## Sumber resmi Meta (wajib dibaca operator)

| Topik | URL resmi |
|-------|-----------|
| Customer Information Parameters | https://developers.facebook.com/docs/marketing-api/conversions-api/parameters/customer-information-parameters |
| Server Event Parameters | https://developers.facebook.com/docs/marketing-api/conversions-api/parameters/server-event |
| Semua Parameters (indeks) | https://developers.facebook.com/docs/marketing-api/conversions-api/parameters |
| Best Practices (wajib + rekomendasi) | https://developers.facebook.com/docs/marketing-api/conversions-api/best-practices |
| fbp & fbc | https://developers.facebook.com/docs/marketing-api/conversions-api/parameters/fbp-and-fbc |
| Dedup Pixel + CAPI | https://developers.facebook.com/docs/marketing-api/conversions-api/deduplicate-pixel-and-server-events |
| Parameter Builder (validasi JSON) | https://developers.facebook.com/docs/marketing-api/conversions-api/parameters/parameter-builder |
| Dataset Quality API | https://developers.facebook.com/docs/marketing-api/conversions-api/dataset-quality-api |

Dokumentasi baru juga tersedia di path `developers.facebook.com/documentation/ads-commerce/conversions-api/...` (struktur 2025+); isi setara dengan path `/docs/marketing-api/...` di atas.

---

## 1. Kesimpulan: IP + Device Saja = Tidak Sesuai Meta

| Yang sering salah | Kebenaran Meta (2026) |
|-------------------|------------------------|
| Cukup `client_ip_address` + `client_user_agent` | **Tidak cukup** untuk EMQ — hanya pelengkap |
| Semua field di-hash | **Salah** — IP, UA, `fbp`, `fbc` **tidak** di-hash |
| Email dikirim plain | **Salah** — `em` wajib SHA256 setelah normalisasi |
| `user_data` boleh kosong | **Salah** — minimal **satu** parameter customer info dengan format benar |
| Pixel browser tidak perlu | Untuk website, Meta wajibkan juga `event_source_url` + `client_user_agent` di CAPI |

**Pixel Hub Seosementara** harus mengumpulkan dan mengirim setara **Advanced Matching + CAPI Best Practices**, bukan forwarder IP.

---

## 2. Parameter Wajib — Event Website (Server / CAPI)

Menurut [Parameters](https://developers.facebook.com/docs/marketing-api/conversions-api/parameters) dan [Best Practices](https://developers.facebook.com/docs/marketing-api/conversions-api/best-practices):

### 2.1 Server event (level event, bukan hanya `user_data`)

| Parameter | Wajib website? | Hash? | Keterangan |
|-----------|----------------|-------|------------|
| `event_name` | Ya | - | Standard atau custom |
| `event_time` | Ya | - | Unix detik, waktu aksi user |
| `action_source` | Ya | - | Website = `website` |
| `event_source_url` | **Ya (website)** | - | URL halaman, harus selaras domain terverifikasi |
| `client_user_agent` | **Ya (website)** | **Tidak** | Di dalam objek `user_data` |
| `event_id` | Sangat disarankan | - | Dedup dengan Meta Pixel |
| `event_source_url` | Ya | - | |

### 2.2 Customer information (`user_data`)

| Aturan Meta | Detail |
|-------------|--------|
| Minimal satu parameter | Harus salah satu dari daftar §3 dengan **format benar** |
| Graph API v13+ | Ada aturan **kombinasi** parameter yang dianggap valid — ikuti Best Practices |
| Contact info | `em`, `ph`, `fn`, `ln`, `ct`, `st`, `zp`, `country`, `ge`, `db` → **wajib hash** |
| Teknis / cookie | `client_ip_address`, `client_user_agent`, `fbp`, `fbc` → **jangan hash** |
| `external_id` | Hash **disarankan** (bukan wajib hash di semua kasus, tapi praktik Pro: hash) |

Meta juga merekomendasikan **`external_id` + `event_id`** untuk semua event.

---

## 3. Tabel Lengkap `user_data` (Sesuai Meta 2026)

| Key API | Label | Hash? | Normalisasi sebelum hash | Prioritas EMQ |
|---------|-------|-------|-------------------------|---------------|
| `em` | Email | **Ya** | Trim, lowercase | **Tertinggi** |
| `ph` | Telepon | **Ya** | Hanya digit + **kode negara** (contoh US: `1` + nomor) | **Tertinggi** |
| `fn` | Nama depan | **Ya** | Lowercase, tanpa tanda baca, UTF-8 | Sedang |
| `ln` | Nama belakang | **Ya** | Sama `fn` | Sedang |
| `ge` | Gender | **Ya** | `m` / `f` lowercase | Sedang |
| `db` | Tanggal lahir | **Ya** | Format `YYYYMMDD` | Sedang |
| `ct` | Kota | **Ya** | Lowercase, tanpa spasi berlebih | Sedang |
| `st` | Provinsi | **Ya** | Kode 2 huruf (US) | Sedang |
| `zp` | Kode pos | **Ya** | Min 5 digit (US) | Sedang |
| `country` | Negara | **Ya** | ISO 3166-1 alpha-2 **lowercase** (`id`, `us`) | Sedang |
| `external_id` | ID pengguna CMS | Disarankan hash | ID stabil (user_id, member_id) | Tinggi |
| `client_ip_address` | IP | **Tidak** | IPv6 lebih disarankan jika ada | Sedang |
| `client_user_agent` | UA | **Tidak** | String lengkap | Sedang |
| `fbp` | Browser ID Meta | **Tidak** | Cookie `_fbp` — **refresh** berkala | **Tinggi** |
| `fbc` | Click ID | **Tidak** | Dari `_fbc` atau `fbclid` | **Tinggi** (iklan) |
| `subscription_id` | ID langganan | **Tidak** | Untuk event subscribe | Khusus |
| `fb_login_id` | Facebook Login ID | **Tidak** | Jika pakai Login Facebook | Khusus |
| `lead_id` | Lead ID | **Tidak** | **Lead Ads / CRM CAPI** | Wajib untuk lead optimization |
| `page_id` | Page ID | **Tidak** | Messaging / page scope | Khusus |
| `page_scoped_user_id` | PSID | **Tidak** | Messenger / page | Khusus |
| `ctwa_clid` | Click to WhatsApp | **Tidak** | Iklan WA | Khusus |
| `ig_account_id` / `ig_sid` | Instagram | **Tidak** | Iklan IG messaging | Khusus |

**Catatan:** `madid`, `anon_id` hanya untuk **app events**, bukan website biasa.

---

## 4. Normalisasi & Hash (Implementasi Hub — Wajib Benar)

Salah normalisasi = hash beda = Meta **tidak match** (EMQ rendah meskipun “sudah kirim email”).

### 4.1 Email (`em`)

```
Input:  "  User@Example.COM  "
Step 1: trim
Step 2: lowercase → "user@example.com"
Step 3: SHA256 → hex (array di JSON: "em": ["<hex>"])
```

### 4.2 Telepon (`ph`) — Indonesia contoh

```
Input:  "+62 812-3456-7890" / "081234567890"
Step 1: buang non-digit kecuali leading country
Step 2: hasil digit saja dengan kode negara: "6281234567890"
Step 3: SHA256 → "ph": ["<hex>"]
```

Meta: **selalu** sertakan kode negara, meskipun semua user dari satu negara.

### 4.3 Nama (`fn` / `ln`)

- Lowercase, hapus punctuation
- Karakter non-Latin: UTF-8 lalu hash

### 4.4 Yang tidak boleh di-hash

```
client_ip_address  → string plain
client_user_agent  → string plain
fbp                → "fb.1.<creationTime>.<random>"
fbc                → "fb.1.<creationTime>.<fbclid>"
lead_id            → plain (dari Lead Ads)
```

### 4.5 `fbp` / `fbc` (Meta first-party cookie)

| Cookie | Format | Cara dapat |
|--------|--------|------------|
| `_fbp` | `fb.1.{unix}.{random}` | Set oleh Meta Pixel / Parameter Builder / first-party |
| `_fbc` | `fb.1.{unix}.{click_id}` | Dari parameter URL `fbclid` pada landing iklan |

Meta: nilai **berubah** — harus **di-refresh** ke CAPI, bukan cache sekali di server selamanya.

---

## 5. `custom_data` — Standard Events (Purchase, Lead, dll.)

Selain `user_data`, Meta memakai `custom_data` untuk optimasi nilai & katalog.

| Event | Field `custom_data` penting |
|-------|----------------------------|
| `Purchase` | `value`, `currency`, `order_id`, `content_ids`, `num_items` |
| `Lead` | `content_name`, `content_category` (opsional) |
| `ViewContent` | `content_ids`, `content_type`, `content_name` |
| `AddToCart` | `content_ids`, `value`, `currency` |
| `InitiateCheckout` | `value`, `num_items` |
| `Search` | `search_string` |

`currency`: ISO 4217 (`IDR`, `USD`). `value`: angka (float).

---

## 6. Dedup Resmi (Browser Pixel + CAPI)

[Best Practices](https://developers.facebook.com/docs/marketing-api/conversions-api/best-practices):

| Syarat | Detail |
|--------|--------|
| `event_name` | **Identik** antara browser dan server |
| Dedup key | **`event_id`** **atau** kombinasi **`external_id` + `fbp`** |
| Rekomendasi Meta | Kirim **`event_id` + `external_id` + `fbp`** sekaligus |

Pixel Hub:

1. Generate `event_id` di ingest  
2. Pass ke browser (`eventID`) jika hybrid  
3. Kirim CAPI dengan `event_id` sama  
4. Sertakan `external_id` (user login) + `fbp` dari cookie  

---

## 7. Tier Kualitas — Selaras Meta (bukan estimasi sembarangan)

| Tier | Isi (minimum) | Layak optimasi iklan? |
|------|---------------|------------------------|
| **D — Invalid untuk Meta Pro** | Hanya IP + UA, tanpa minimal 1 customer param valid | **Tidak** |
| **C — Minimum valid** | + `event_source_url` + `action_source` + (`fbp` **atau** hash `em`) | Testing / traffic |
| **B — Recommended** | + `fbc` (jika ads) + hash `em`/`ph` + `event_id` + `external_id` | **Lead / retargeting** |
| **A — Optimal** | + `fn`/`ln`/`ct`/`country` + `value`/`order_id` + `lead_id` (lead ads) | **Purchase / scale** |

**Target Seosementara:**

| Jenis event | Tier minimum |
|-------------|--------------|
| `PageView` | C |
| `Lead` | B (+ `lead_id` jika dari Lead Ads) |
| `Purchase` | A |

---

## 8. Spesifikasi Pengumpulan — Pixel Hub (Harus Dibangun)

### 8.1 Browser (`sseo-track.js` + first-party)

| Data | Cara |
|------|------|
| `fbp`, `fbc` | Baca cookie; update dari `fbclid` di URL |
| `event_source_url` | `location.href` |
| URL params | `fbclid`, `gclid` (untuk analytics internal) |
| `session_id` | Cookie first-party Hub |

### 8.2 Server enrich (sebelum CAPI)

| Data | Sumber |
|------|--------|
| `client_ip_address` | `CF-Connecting-IP` / `X-Forwarded-For` |
| `client_user_agent` | Header (wajib website) |
| `em`, `ph` | Form lead/checkout, profil user |
| `fn`, `ln`, `ct`, `zp`, `country` | Form (opsional, naikkan EMQ) |
| `external_id` | `users.id` / member_id |
| `lead_id` | Webhook Lead Ads / CRM [Conversion Leads](https://developers.facebook.com/docs/marketing-api/conversions-api/conversion-leads-integration) |
| `fb_login_id` | Jika integrasi Facebook Login |
| `custom_data` | Order API, produk, shortlink metadata |

### 8.3 Tabel `pixel_sessions` (refresh fbp/fbc)

Simpan per `session_id`: `fbp`, `fbc`, `last_url`, `first_fbclid` — TTL 7–90 hari.

### 8.4 Privacy gateway (sebelum hash)

| Input mentah | Output ke CAPI |
|--------------|----------------|
| email | `em` hashed |
| phone | `ph` hashed |
| IP, UA, fbp, fbc | plain |
| Tanpa consent (GDPR) | Jangan kirim `em`/`ph` — event `skipped` |

---

## 9. Payload Contoh — Selaras Meta (Purchase website)

```json
{
  "data": [{
    "event_name": "Purchase",
    "event_time": 1762902353,
    "event_id": "550e8400-e29b-41d4-a716-446655440000",
    "action_source": "website",
    "event_source_url": "https://rezekibelanja.com/checkout/success",
    "user_data": {
      "em": ["7b17fb0bd173f625b58625ad059fcbc2e2c25691cddad1961d840fcffd356b98"],
      "ph": ["c051715cc583c6386f63ae2e614361fd9a67efb3a9bafa64e36f97c9b4da82c9"],
      "fn": ["51b03d7eafc121fea0e80a5ea83beb7c449f4ec"],
      "client_ip_address": "203.0.113.10",
      "client_user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
      "fbp": "fb.1.1762902000.1987654321",
      "fbc": "fb.1.1762901800.IwAR2xxxxxxxx",
      "external_id": ["a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"],
      "country": ["0425f1493e1614eb8b8f6b7d8e6e5a2e8c8c8c8c8c8c8c8c8c8c8c8c8c8c8"]
    },
    "custom_data": {
      "currency": "IDR",
      "value": 150000,
      "order_id": "ORD-2026-9988",
      "content_ids": ["SKU-1", "SKU-2"],
      "num_items": 2
    }
  }]
}
```

*(Nilai hash di atas ilustrasi — gunakan normalisasi §4 di production.)*

---

## 10. Validasi & Monitoring (Meta Tools)

| Tool | Fungsi di operasi |
|------|-------------------|
| **Test Events** (Events Manager) | `test_event_code` di root payload |
| **Payload Helper / Parameter Builder** | Validasi struktur JSON sebelum production |
| **Dataset Quality API** | Skor kualitas dataset / diagnostics programmatic |
| **Events Manager Diagnostics** | EMQ per parameter — bandingkan dengan §7 tier |

Di admin Hub tab **Connection**: tampilkan checklist parameter Meta (✓ `em` ✓ `fbp` ✗ `ph`) — mirror Events Manager.

---

## 11. Gap: Dokumen / Rencana Lama vs Meta 2026

| Topik | Sebelumnya di Plan kita | Perbaikan (dokumen ini) |
|-------|-------------------------|-------------------------|
| Field `user_data` | Sebagian (`em`, `ph`, `fbp`) | **Tabel lengkap §3** + messaging/lead |
| Hash IP | Kadang disalah pahami | Eksplisit **jangan hash** §4.5 |
| Wajib website | Tidak tegas | `event_source_url` + `client_user_agent` §2 |
| `lead_id` | Tidak ada | Lead Ads + CRM §3 |
| `fb_login_id` | Tidak ada | §3 |
| Dedup | Hanya `event_id` | + `external_id` + `fbp` §6 |
| Graph API v13 kombinasi | Tidak disebut | §2.2 |
| Parameter Builder | Tidak ada | §10 |
| Refresh fbp/fbc | Disebut singkat | §4.5 wajib |

---

## 12. Checklist Implementasi Hub (Wajib untuk “Sesuai Meta”)

### Ingest & browser

- [ ] Kirim `event_source_url`, `client_user_agent` pada **semua** event website
- [ ] Baca `_fbp`, `_fbc`; bangun `fbc` dari `fbclid`
- [ ] Refresh cookie ke session store
- [ ] `event_id` UUID setiap event

### Privacy & hash

- [ ] Normalisasi §4 sebelum SHA256 (`em`, `ph`, `fn`, `ln`, …)
- [ ] **Jangan** hash IP, UA, fbp, fbc
- [ ] Kode negara telepon Indonesia `62` konsisten

### Konversi

- [ ] Form lead → `em` (minimal)
- [ ] Checkout → `em`/`ph` + `custom_data.value/currency/order_id`
- [ ] Lead Ads → `lead_id` + [Conversion Leads](https://developers.facebook.com/docs/marketing-api/conversions-api/conversion-leads-integration)

### Admin & QA

- [ ] Coverage % per parameter (target §7)
- [ ] Alert tier D > 10%
- [ ] Uji di Test Events + Parameter Builder

---

## 13. Dokumen terkait

- [23-meta-conversions-api-kedalaman.md](./23-meta-conversions-api-kedalaman.md)
- [21-pixel-facebook-pro.md](./21-pixel-facebook-pro.md)
- [24-meta-akun-bm-pixel-dan-optimasi-iklan.md](./24-meta-akun-bm-pixel-dan-optimasi-iklan.md)

**Versi dokumen:** 2.0 (selaras Meta Parameters & Best Practices, verifikasi web Mei 2026)
