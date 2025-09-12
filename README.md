# Pryx

Pryx is a lightweight gateway for managing and proxying LLM model requests.
It adds **user management**, **API keys**, and **scope-based authorization** on top of simple HTTP endpoints.

## Features

* REST endpoints for completions and model registration.
* API key authentication (hashed in Postgres).
* Scopes for granular access control (`completion:invoke`, `model:write`, `model:<name>`).
* Automatic database migrations via `pryx-migrate`.
* Modular Go codebase with Chi + GORM.

## Quickstart

```bash
docker-compose up --build
```

## Admin bootstrap

Create a user:

```bash
curl -s http://localhost:8080/admin/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"me@example.com","name":"Me"}'
```

Create an API key:

```bash
curl -s http://localhost:8080/admin/keys \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"name":"cli","scopes":"completion:invoke,model:write"}'
```

Call the service:

```bash
curl -s http://localhost:8080/ \
  -H "Authorization: Bearer sk_live_..." \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hello"}]}'
```

## Documentation

* [Architecture](docs/ARCHITECTURE.md)
* [Authentication](docs/AUTHENTICATION.md)
* [Database](docs/DATABASE.md)
* [Usage](docs/USAGE.md)
