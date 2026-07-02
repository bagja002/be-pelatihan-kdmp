# KNMP Backend — Go Fiber Clean Architecture Template

A Go + [Fiber](https://gofiber.io) starter using **clean architecture** with a
code generator. Adding a new entity scaffolds its repository, service, handler,
DTO and grouped route automatically — **you only register the route.**

## Stack

- **Fiber v2** — HTTP framework
- **GORM + MySQL** — ORM & database
- **go-playground/validator** — request validation
- **godotenv** — env config

## Layout

```
cmd/
  api/          # server entrypoint
  generate/     # entity code generator (+ templates/)
internal/
  config/       # env config
  database/     # GORM connection + auto-migration registry
  entity/       # domain models (self-register for migration)
  dto/          # request/response payloads + validation tags
  repository/   # data access (interface + GORM impl)
  service/      # business logic (interface + impl)
  handler/      # Fiber HTTP handlers
  router/       # route groups (one file per entity) + router.go
pkg/
  response/     # standard JSON envelope
  validator/    # validator wrapper
```

Flow: `handler → service → repository → entity`. Each layer is an interface,
so layers are decoupled and testable.

## Getting started

```bash
cd backend
cp .env.example .env      # then edit DB credentials
make keygen               # copy JWT_SECRET + ENCRYPTION_KEY into .env
make tidy                 # go mod tidy (download deps)
make run                  # or: go run ./cmd/api
```

Server runs on `:3000` (configurable via `APP_PORT`). Health check: `GET /health`.

> In `APP_ENV=production` the app refuses to start unless `JWT_SECRET`
> (≥ 32 chars) and `ENCRYPTION_KEY` (32 bytes) are set. In development it
> falls back to insecure defaults and prints a loud warning.

## Adding a new entity

1. Generate the slice:

   ```bash
   make gen name=Category
   # or: go run ./cmd/generate -name Category
   # custom plural: go run ./cmd/generate -name Category -plural categories
   ```

   This creates:
   - `internal/entity/category.go`
   - `internal/dto/category_dto.go`
   - `internal/repository/category_repository.go`
   - `internal/service/category_service.go`
   - `internal/handler/category_handler.go`
   - `internal/router/category_route.go`

2. **Register the route (the only manual step)** — add one line inside
   `SetupRoutes` in `internal/router/router.go`:

   ```go
   RegisterCategoryRoutes(api, db)
   ```

3. That's it. Migration is automatic (the entity self-registers via `init()`),
   and you immediately get a full CRUD group under `/api/v1/categories`.

> Adjust the generated `entity`, `dto`, and the field mapping in the `service`
> to add real fields — the generator scaffolds a single `Name` field as a
> starting point.

## Endpoints (per entity, e.g. `products`)

| Method | Path                     | Action        |
| ------ | ------------------------ | ------------- |
| POST   | `/api/v1/products`       | create        |
| GET    | `/api/v1/products`       | list          |
| GET    | `/api/v1/products/:id`   | detail        |
| PUT    | `/api/v1/products/:id`   | update        |
| DELETE | `/api/v1/products/:id`   | delete        |

All responses use a consistent envelope:

```json
{ "success": true, "message": "product detail", "data": { } }
```

## Security

Application-layer hardening built in (aligned with OWASP Top 10 & modern web standards):

| Area | Implementation |
| ---- | -------------- |
| **Authentication** | JWT (HS256) access + refresh tokens, token-type enforcement, alg-confusion protection (`pkg/token`) |
| **Password storage** | Argon2id (OWASP-recommended), per-password random salt, constant-time verify (`pkg/hash`) |
| **Data-at-rest encryption** | AES-256-GCM field-level encryption via GORM serializer; tag a field `gorm:"serializer:encrypted"` (`pkg/crypto`, `internal/database/serializer.go`) |
| **Security headers** | Helmet: HSTS (preload), CSP, `X-Frame-Options`, `X-Content-Type-Options: nosniff`, `Referrer-Policy` |
| **CORS** | Explicit origin allow-list (never wildcard with credentials) |
| **Rate limiting** | Global per-IP limiter + stricter limiter on `/auth/*` (anti brute-force) |
| **Info-leak prevention** | Global error handler + `response.InternalError` log internally, return generic messages; internal errors never sent to clients |
| **Audit logging** | Per-request `X-Request-ID` propagated into access + error logs |
| **DoS mitigation** | Request body size limit + read/write/idle timeouts |
| **Secret hygiene** | Secrets from env; production fails fast on missing/weak secrets; password hash & DB password never serialized (`json:"-"`) |
| **SQL injection** | GORM parameterizes all queries |

### Auth endpoints

| Method | Path                     | Auth   | Action                      |
| ------ | ------------------------ | ------ | --------------------------- |
| POST   | `/api/v1/auth/register`  | public | create account              |
| POST   | `/api/v1/auth/login`     | public | get access + refresh tokens |
| POST   | `/api/v1/auth/refresh`   | public | rotate tokens               |
| GET    | `/api/v1/auth/me`        | Bearer | current user profile        |

Send the access token as `Authorization: Bearer <token>` on protected routes.

### Protecting entity routes

In `internal/router/router.go` there are two groups — register your entity on
whichever fits:

```go
RegisterProductRoutes(protected, db) // requires a valid JWT
RegisterProductRoutes(api, db)       // public
```

### Encrypting a sensitive field

Add the serializer tag to any string field (see `internal/entity/user.go` `Phone`):

```go
SSN string `gorm:"type:varchar(512);serializer:encrypted" json:"ssn,omitempty"`
```

It is stored as AES-256-GCM ciphertext and transparently decrypted on read.

## Make targets

| Command                | Description                     |
| ---------------------- | ------------------------------- |
| `make run`             | run the API server              |
| `make gen name=Foo`    | scaffold a new entity           |
| `make tidy`            | resolve go module dependencies  |
| `make build`           | build binary into `./bin/api`   |
| `make keygen`          | print fresh JWT + encryption secrets |
