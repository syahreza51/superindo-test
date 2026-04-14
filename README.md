# superindo-test

API Golang untuk endpoint `/product` (create + list + search + filter + sort) dengan:
- Database: PostgreSQL (migrations + seeder)
- Cache: Redis
- DI: wire
- Unit test
- Docker

## Menjalankan dengan Docker

1) Salin env:

```bash
cp .env.example .env
```

2) Jalankan:

```bash
docker compose up --build
```

Services:
- `api`: `http://localhost:8080`
- `postgres`: `localhost:5432`
- `redis`: `localhost:6379`

## API

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

### GET /product
Ambil list produk dengan query opsional:
- `id`: exact match by UUID
- `name`: partial match (ILIKE)
- `type`: `Sayuran|Protein|Buah|Snack` (bisa diulang: `type=Buah&type=Snack`)
- `sort_by`: `date|price|name` (default `date`)
- `sort_order`: `asc|desc` (default `desc`)
- `limit`: default 20, max 100
- `offset`: default 0

Contoh:
- `GET /product?name=ayam&type=Protein&sort_by=price&sort_order=asc`

# superindo-test