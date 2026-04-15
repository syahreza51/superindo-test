# superindo-test

API Golang untuk endpoint `POST /product` dan `GET /product` (create + list + search + filter + sort) dengan PostgreSQL (migrations + seed) dan Redis (cache).

## Fitur

- **REST API**: create & list product
- **Search + filter + sort**: via query parameter di `GET /product`
- **PostgreSQL**: migrations + seeder
- **Redis cache**: configurable TTL
- **Dependency Injection**: `google/wire`
- **Unit test**: ada test untuk modul product
- **Docker Compose**: ready-to-run (db + redis + migrate + seed + api)

## Tech stack

- **Go**: 1.23
- **HTTP router**: `chi`
- **Database**: PostgreSQL 16
- **Cache**: Redis 7

## Menjalankan dengan Docker (recommended)

1) Salin env:

```bash
cp .env.example .env
```

2) Jalankan:

```bash
docker compose up --build
```

Docker compose akan menjalankan urutan berikut:
- **postgres** → siap (healthcheck)
- **migrate** → apply migrations dari folder `migrations/`
- **seed** → insert data dari `seed/seed.sql`
- **redis** → siap (healthcheck)
- **api** → start server

Service yang tersedia:
- **api**: `http://localhost:8080` (atau sesuai `APP_PORT`)
- **postgres**: `localhost:5432`
- **redis**: `localhost:6379`

## Menjalankan tanpa Docker (opsional)

Pastikan kamu punya:
- **Go 1.23**
- **PostgreSQL** dan **Redis** yang bisa diakses dari app

1) Buat `.env`:

```bash
cp .env.example .env
```

2) Sesuaikan `DB_HOST` dan `REDIS_ADDR` agar mengarah ke service lokal (contoh: `localhost:5432` dan `localhost:6379`).

3) Jalankan API:

```bash
go run ./cmd/api
```

Catatan: migrations & seed pada repo ini sudah disiapkan untuk flow Docker Compose (via container `migrate` dan `psql`). Kalau run tanpa Docker, jalankan migrations/seed sesuai tooling kamu.

## Environment variables

Contoh lengkap ada di `.env.example`. Variabel yang paling sering diubah:
- **APP_PORT**: port yang di-expose (default `8080`)
- **DB_\***: koneksi PostgreSQL
- **REDIS_ADDR**: address Redis (default compose: `redis:6379`)
- **REDIS_TTL_SECONDS**: TTL cache (default `30`)

## Project structure

Struktur folder utama:

```text
cmd/api/          entrypoint aplikasi (main + wire)
internal/app/     wiring & bootstrap aplikasi (router, config, dependencies)
internal/product/ domain product (handler, service, repository, cache, tests)
migrations/       SQL migrations untuk PostgreSQL
seed/             SQL seed data
```

## API

Base URL (Docker): `http://localhost:8080`

### POST /product

Tambah produk.

Body:

```json
{
  "name": "Bayam",
  "type": "Sayuran",
  "price": 12000
}
```

Contoh `curl`:

```bash
curl -sS -X POST "http://localhost:8080/product" \
  -H "Content-Type: application/json" \
  -d '{"name":"Bayam","type":"Sayuran","price":12000}'
```

### GET /product

Ambil list produk dengan query opsional:
- **id**: exact match by UUID
- **name**: partial match (ILIKE)
- **type**: `Sayuran|Protein|Buah|Snack` (bisa diulang: `type=Buah&type=Snack`)
- **sort_by**: `date|price|name` (default `date`)
- **sort_order**: `asc|desc` (default `desc`)
- **limit**: default 20, max 100
- **offset**: default 0

Contoh:

```bash
curl -sS "http://localhost:8080/product?name=ayam&type=Protein&sort_by=price&sort_order=asc"
```

## Testing

```bash
go test ./...
```