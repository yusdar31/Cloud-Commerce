# User Service (Identity Service)

**Service:** user-service

**Port:** `8081`

**Database:** `identity_db` (PostgreSQL)

**Owner:** Backend Team

---

## Responsibility

Mengelola autentikasi, otorisasi, dan manajemen pengguna.

- Register pengguna baru (Seller / Buyer)
- Login dengan email & password
- Generate JWT access token + refresh token
- Refresh token rotation
- Logout (blacklist token di Redis)
- Manajemen profil pengguna
- Role management: `seller`, `buyer`, `admin`

---

## API Endpoints

Base URL: `http://localhost:8081`

### Public (Tidak butuh auth)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/register` | Register pengguna baru |
| `POST` | `/api/v1/auth/login` | Login, dapatkan JWT |
| `POST` | `/api/v1/auth/refresh` | Refresh access token |

### Protected (Butuh JWT)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/auth/logout` | Logout, blacklist token |
| `GET` | `/api/v1/users/me` | Profil pengguna saat ini |
| `PUT` | `/api/v1/users/me` | Update profil |

### Health

| Method | Path |
|--------|------|
| `GET` | `/healthz` |

---

## Environment Variables

| Variable | Wajib | Default (Dev) |
|----------|-------|---------------|
| `PORT` | ✅ | `8081` |
| `DATABASE_URL` | ✅ | `postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable` |
| `REDIS_URL` | ✅ | `redis://localhost:6379` |
| `JWT_SECRET` | ✅ | `dev-secret-change-me` |
| `JWT_EXPIRY` | ✅ | `24h` |
| `REFRESH_TOKEN_EXPIRY` | ❌ | `720h` |
| `APP_ENV` | ✅ | `development` |

---

## Running Locally

```bash
cd services/user-service

export PORT=8081
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/identity_db?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="dev-secret-change-me"
export JWT_EXPIRY="24h"
export APP_ENV="development"

goose -dir migrations postgres "$DATABASE_URL" up
go run ./cmd/server
```

---

## JWT Claims Structure

```json
{
  "sub": "usr_01J6...",
  "tenant_id": "ten_01J6...",
  "role": "seller",
  "email": "seller@example.com",
  "iat": 1720000000,
  "exp": 1720086400
}
```

---

## Related Documents

- [Backend Guidelines](../../docs/04-engineering/04-backend-guidelines.md)
- [API Guidelines](../../docs/03-architecture/api-guidelines.md)
- [Environment Variables](../../docs/04-engineering/env-variables-reference.md)
