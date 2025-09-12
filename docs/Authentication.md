# Authentication & Authorization

Pryx uses API keys to secure endpoints.

## API Key Format

* Keys are generated with a `sk_live_` prefix followed by 64 hex characters.
* Example: `sk_live_abc123...`
* Keys are only shown **once** at creation.

## Storage

* Stored in DB as:

  * `prefix`: first 8 characters (for lookup/logging).
  * `hash`: SHA-256 hash of the full key (no plaintext).
  * `scopes`: comma-separated list of permissions.
  * `revoked`: boolean to disable keys.
  * `last_used_at`: timestamp updated on request.

## Middleware

Each protected route attaches `auth.Middleware(db, requireScope)`:

* Checks bearer token from `Authorization: Bearer <token>`.
* Validates key prefix and hash against DB.
* Verifies required scope if specified.
* Injects `AuthContext` into `context.Context` for downstream handlers.

## Scopes

* `completion:invoke` → call `/` for completions.
* `model:write` → register new models.
* `model:any` → call any model.
* `model:<name>` → call only specific model.

## Example Usage

```http
Authorization: Bearer sk_live_abcd1234...
```

If the key lacks the required scope, server responds with:

```json
{"error":"forbidden: missing scope"}
```

## Revocation

Keys can be revoked by setting `revoked=true` in the database. Middleware blocks revoked keys immediately.

## Best Practices

* Treat keys like passwords. Never log the full key.
* Use scopes to enforce least privilege.
* Rotate keys regularly and revoke old ones.
* Use `prefix` for identification in logs and UIs.
