# Seosementara

CMS untuk operasi massal domain & iklan. Dokumentasi: `Plan/`. Kode aktif: **Pixel Facebook Pro** (MVP).

## Mulai cepat — Pixel Facebook Pro

```bash
cd Backend
openssl rand -base64 32   # simpan sebagai PIXEL_ENCRYPTION_KEY
export PIXEL_ENCRYPTION_KEY="..."
export ADMIN_TEMPLATES_DIR="../Frontend-admin/templates"
export STATIC_DIR="../Frontend-admin/static"
go run ./cmd/api
```

Admin: http://localhost:8080/admin/pixel/facebook/