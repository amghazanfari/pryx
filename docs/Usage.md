# Usage Guide

This guide explains how to run Pryx, bootstrap users and API keys, and interact with its HTTP endpoints.

## Start the stack

```bash
docker-compose up --build
```

Services:

* **pryx**: API server on `:8080`
* **pryx-migrate**: runs schema migrations once
* **postgres**: database

Optional dev setting:

* `DB_AUTOMIGRATE=true` to let the app run `AutoMigrateAll()` at startup.

## Environment variables

Provided via `docker-compose.yml` or your shell:

* `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
* `POSTGRES_SSLMODE` (default `disable`)
* `DB_TIMEZONE` (default `UTC`)

## Admin bootstrap

Create a user:

```bash
curl -s http://localhost:8080/admin/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"me@example.com","name":"Me"}'
```

Create an API key (save the plaintext once):

```bash
curl -s http://localhost:8080/admin/keys \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"name":"cli","scopes":"completion:invoke,model:any,model:write"}'
```

Response includes `{ id, prefix, key }`. Store `key` securely; it is not retrievable later.

## Authentication header

All protected endpoints require:

```http
Authorization: Bearer <api_key>
```

## Endpoints

### 1) POST `/v1/completions` – Chat completions

* **Scope required**: `completion:invoke`
* **Body**: compatible with OpenAI `chat.completions.create` style payload

Example:

```bash
curl -s http://localhost:8080/v1/completions \
  -H "Authorization: Bearer sk_live_..." \
  -H "Content-Type: application/json" \
  -d '{
    "model":"gpt-4o-mini",
    "messages":[{"role":"user","content":"Hello!"}]
  }'
```

Responses are JSON pass-through from the configured upstream.

Per-model authorization (if enabled): the handler inspects `model` and requires `model:<name>` or `model:any` scope.

### 2) POST `/v1/models` – Register an upstream model

* **Scope required**: `model:write`
* **Body**:

```json
{
  "name": "OpenRouter GPT-4o Mini",
  "model_name": "gpt-4o-mini",
  "endpoint": "https://openrouter.ai/api/v1",
  "api_key": ""  
}
```

Example:

```bash
curl -s http://localhost:8080/v1/models \
  -H "Authorization: Bearer sk_live_..." \
  -H "Content-Type: application/json" \
  -d '{"name":"OpenRouter GPT-4o Mini","model_name":"gpt-4o-mini","endpoint":"https://openrouter.ai/api/v1","api_key":""}'
```

> Note: if your deployment uses a single upstream configured by env, `/models` can still serve as a registry/audit.

## Error responses

* Missing or malformed token → `401 Unauthorized`
* Valid token without scope → `403 Forbidden`
* Bad JSON payload → `400 Bad Request`
* Unknown model (when enforced) → `400 Bad Request`
* Upstream failure → `502 Bad Gateway`

## Scope reference

| Scope               | Grants                            |
| ------------------- | --------------------------------- |
| `completion:invoke` | Call the `/` completions endpoint |
| `model:write`       | Register models via `/models`     |
| `model:any`         | Use any model name                |
| `model:<name>`      | Use only the specified model      |

## Revoking and rotating keys

* Revoke: set `revoked=true` for the key in the database.
* Rotate: create a new key, update clients, revoke the old key.

## Health and logs

* Pryx logs JSON to stdout (via logrus).
* Look for request logs and 4xx/5xx status codes.
* You can add your own `/healthz` route in `cmd/app/main.go` if needed.

## Troubleshooting

* **401**: missing/invalid token, or key revoked. Check `Authorization` header and DB.
* **403**: scope missing. Ensure key includes the required scope(s).
* **502**: upstream error. Verify upstream `endpoint`, network, and credentials.
* **Model not allowed**: add `model:any` or `model:<name>` to the key.

## Production recommendations

* Run migrations via `pryx-migrate` in CI/CD or init jobs.
* Put Pryx behind a reverse proxy with TLS.
* Do not expose `/admin/*` publicly; gate with network policy or a shared secret.
* Add rate limiting per `api_key_id` (e.g., Redis token bucket).
* Track usage by persisting token counts and latency per request.
