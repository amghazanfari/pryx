# Architecture

Pryx is organized around a **Go service + Postgres backend**.

## Components

* **cmd/app**
  Entry point for the API server. Wires routes and middleware.

* **cmd/migrate**
  Applies migrations using `gormigrate`.

* **internal/db**
  Database connection handling and auto-migrations.

* **internal/handlers**
  HTTP endpoints (`CompletionHandler`, `AddModelHandler`, and admin endpoints for users and keys).

* **internal/auth**
  API key generation, hashing, scope checks, Chi middleware for authentication and authorization.

* **internal/models**
  Database models: `User`, `APIKey`, `Model`.

## Request Flow

1. Client sends HTTP request with `Authorization: Bearer <api_key>`.
2. `auth.Middleware` validates the key and checks required scopes.
3. Handler executes business logic (model lookup, forwarding to upstream completion endpoint, or registering models).
4. Response returned as JSON.

## Security Layers

* **API Keys**: scoped permissions (`completion:invoke`, `model:write`, `model:any`, `model:<name>`).
* **Hashing**: Keys stored hashed with SHA-256, only prefix stored for lookup.
* **Middleware**: Enforces authentication before business logic.

## Deployment

* **docker-compose.yml** orchestrates:

  * `pryx`: main API server (port 8080).
  * `pryx-migrate`: migration runner.
  * `postgres`: persistence layer.

* **Scaling**: `pryx` can be scaled horizontally behind a load balancer. Postgres should run with proper backups and connection pooling.

## Extensibility

* Add new scopes for fine-grained control.
* Add new migrations with `gormigrate`.
* Extend `internal/handlers` for additional endpoints (usage metrics, revocation, etc.).
