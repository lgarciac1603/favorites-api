# Favorites API

A production-ready REST API built with Go, Gin, and PostgreSQL that manages a user's favorite cryptocurrencies. This service operates as an **optional microservice** alongside [cpp-rest-api](https://github.com/lgarciac1603/cpp-rest-api), the primary authentication and user management backend.

> **Optional microservice**: `favorites-api` is not required for `cpp-rest-api` to function. You can run `cpp-rest-api` independently and add this microservice when you need favorites management. When running standalone, `favorites-api` requires `cpp-rest-api` to be reachable at `:8080` for token validation. When running as part of the full stack (via `crypto-dashboard`), both services are orchestrated automatically.

---

## Table of Contents

1. [Overview](#overview)
2. [Technology Stack](#technology-stack)
3. [Prerequisites](#prerequisites)
4. [Installation](#installation)
5. [Configuration](#configuration)
6. [Project Structure](#project-structure)
7. [API Endpoints](#api-endpoints)
8. [Running the Application](#running-the-application)
9. [Testing](#testing)
10. [Integration with cpp-rest-api](#integration-with-cpp-rest-api)
11. [Database Schema](#database-schema)
12. [Authentication Flow](#authentication-flow)
13. [Best Practices](#best-practices)
14. [Contributing](#contributing)
15. [License](#license)
16. [Roadmap](#roadmap)

---

## Overview

Favorites API is an **optional auxiliary microservice** that complements **cpp-rest-api**, a C++ REST service that owns user management and authentication. Its sole responsibility is managing the `user_favorites` table in a shared PostgreSQL instance. It does not handle user creation, login, or token issuance — those concerns remain entirely in `cpp-rest-api`.

`cpp-rest-api` functions fully without this service. `favorites-api` is only needed when cryptocurrency favorites persistence is required.

Core capabilities:

- Retrieve all favorite cryptocurrencies for the authenticated user.
- Add a new cryptocurrency to the user's favorites list.
- Remove a cryptocurrency from the user's favorites list.

The service validates incoming Bearer tokens by delegating to the primary API or, in the current placeholder implementation, by accepting any non-empty token. It then resolves the `userID` from context and scopes every database query to that user.

---

## Technology Stack

| Component         | Technology                                    | Version |
| ----------------- | --------------------------------------------- | ------- |
| Language          | Go                                            | 1.26.1  |
| HTTP Framework    | Gin (`github.com/gin-gonic/gin`)              | 1.12.0  |
| PostgreSQL Driver | lib/pq (`github.com/lib/pq`)                  | 1.12.2  |
| Testing Framework | Testify (`github.com/stretchr/testify`)       | 1.11.1  |
| Database Mocking  | go-sqlmock (`github.com/DATA-DOG/go-sqlmock`) | 1.5.2   |
| Container Runtime | Docker                                        | 26+     |
| Database          | PostgreSQL                                    | 16      |

---

## Prerequisites

Before running this project, ensure the following are installed and available:

- **Go 1.26.1** or later — [https://go.dev/dl/](https://go.dev/dl/)
- **PostgreSQL 16** (or a compatible instance) — [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
- **Docker and Docker Compose** (optional, for containerized runs) — [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)
- **Git** — [https://git-scm.com/downloads](https://git-scm.com/downloads)

---

## Installation

```bash
# 1. Clone the repository
git clone https://github.com/lgarciac1603/favorites-api.git
cd favorites-api

# 2. Download Go module dependencies
go mod download

# 3. (Optional) Verify the module graph
go mod tidy
```

If you intend to run against a local PostgreSQL instance, create the `user_favorites` table before starting the server (see [Database Schema](#database-schema)).

---

## Configuration

The application reads all configuration from environment variables. Defaults are provided for local development.

| Variable   | Default        | Description                     |
| ---------- | -------------- | ------------------------------- |
| `DB_HOST`  | `localhost`    | PostgreSQL host                 |
| `DB_PORT`  | `8090`         | PostgreSQL port                 |
| `DB_NAME`  | `apidb`        | Target database name            |
| `DB_USER`  | `apiuser_test` | Database user                   |
| `DB_PASS`  | `apipass_test` | Database password               |
| `APP_PORT` | `8090`         | HTTP port the server listens on |

For local development, create a `.env` file (not committed to source control) or export each variable in your shell:

```bash
export DB_HOST=localhost
export DB_PORT=8090
export DB_NAME=apidb
export DB_USER=apiuser_test
export DB_PASS=apipass_test
export APP_PORT=8090
```

When running via Docker Compose, these values are injected automatically through the `environment` block in `docker-compose.yml`.

---

## Project Structure

```
favorites-api/
├── config/
│   ├── config.go           # Reads environment variables; builds connection string
│   └── config_test.go      # Unit tests for config package
├── database/
│   ├── database.go         # Opens and manages the PostgreSQL connection
│   └── database_test.go    # Integration tests for database lifecycle
├── docker/
│   └── init.sql            # SQL script that creates user_favorites on first run
├── handlers/
│   ├── favorites.go        # HTTP handlers: GetFavorites, PostFavorite, DeleteFavorite
│   └── favorites_test.go   # Handler unit tests using go-sqlmock
├── middleware/
│   ├── auth.go             # Bearer token validation middleware
│   └── auth_test.go        # Middleware unit tests
├── models/
│   ├── favorite.go         # Favorite struct with JSON tags
│   └── favorite_test.go    # Model serialization tests
├── .dockerignore           # Files excluded from the Docker build context
├── .env.local              # Example local environment file (not committed)
├── .gitignore              # Git ignore rules
├── Dockerfile              # Multi-stage Docker build
├── docker-compose.yml      # Orchestrates postgres + api services
├── go.mod                  # Module definition and dependency versions
├── go.sum                  # Dependency checksums
└── main.go                 # Application entry point: wires config, DB, router
```

---

## API Endpoints

### Base URL

```
http://localhost:8090
```

### Authentication

All endpoints require a valid Bearer token in the `Authorization` header:

```
Authorization: Bearer <token>
```

If the header is absent, malformed, or the token is invalid, the server returns `401 Unauthorized`.

---

### GET /favorites

Retrieves all favorite cryptocurrencies for the authenticated user, ordered by creation date (most recent first).

**Request**

```
GET /favorites
Authorization: Bearer <token>
```

**Response — 200 OK**

```json
{
  "data": [
    {
      "id": 1,
      "userId": 1,
      "cryptoId": "bitcoin",
      "cryptoName": "Bitcoin",
      "createdAt": "2026-04-02T10:30:00Z"
    },
    {
      "id": 2,
      "userId": 1,
      "cryptoId": "ethereum",
      "cryptoName": "Ethereum",
      "createdAt": "2026-04-02T10:35:00Z"
    }
  ],
  "total": 2
}
```

When the user has no favorites, `data` is an empty array and `total` is `0`.

**cURL Example**

```bash
curl -X GET http://localhost:8090/favorites \
  -H "Authorization: Bearer my-secret-token"
```

**JavaScript / TypeScript Example**

```typescript
const response = await fetch("http://localhost:8090/favorites", {
  method: "GET",
  headers: {
    Authorization: "Bearer my-secret-token",
  },
});

const body = await response.json();
console.log(body.data); // Favorite[]
console.log(body.total); // number
```

---

### POST /favorites

Adds a cryptocurrency to the authenticated user's favorites list.

**Request**

```
POST /favorites
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body**

| Field        | Type   | Required | Description                            |
| ------------ | ------ | -------- | -------------------------------------- |
| `cryptoId`   | string | Yes      | Unique identifier (e.g. `"bitcoin"`)   |
| `cryptoName` | string | Yes      | Human-readable name (e.g. `"Bitcoin"`) |

```json
{
  "cryptoId": "ethereum",
  "cryptoName": "Ethereum"
}
```

**Validations**

- Both `cryptoId` and `cryptoName` are required. Missing either field returns `400 Bad Request`.
- A user cannot add the same `cryptoId` twice. A duplicate returns `409 Conflict`.

**Response — 201 Created**

```json
{
  "message": "Crypto added to favorites",
  "data": {
    "id": 3,
    "userId": 1,
    "cryptoId": "ethereum",
    "cryptoName": "Ethereum",
    "createdAt": "2026-04-02T11:00:00Z"
  }
}
```

**cURL Example**

```bash
curl -X POST http://localhost:8090/favorites \
  -H "Authorization: Bearer my-secret-token" \
  -H "Content-Type: application/json" \
  -d '{"cryptoId":"ethereum","cryptoName":"Ethereum"}'
```

**JavaScript / TypeScript Example**

```typescript
const response = await fetch("http://localhost:8090/favorites", {
  method: "POST",
  headers: {
    Authorization: "Bearer my-secret-token",
    "Content-Type": "application/json",
  },
  body: JSON.stringify({ cryptoId: "ethereum", cryptoName: "Ethereum" }),
});

const body = await response.json();
```

---

### DELETE /favorites/{cryptoId}

Removes a cryptocurrency from the authenticated user's favorites list.

**Request**

```
DELETE /favorites/{cryptoId}
Authorization: Bearer <token>
```

**Path Parameter**

| Parameter  | Type   | Description                            |
| ---------- | ------ | -------------------------------------- |
| `cryptoId` | string | The ID of the cryptocurrency to remove |

**Response — 200 OK**

```json
{
  "message": "Crypto removed from favorites"
}
```

**cURL Example**

```bash
curl -X DELETE http://localhost:8090/favorites/ethereum \
  -H "Authorization: Bearer my-secret-token"
```

**JavaScript / TypeScript Example**

```typescript
const response = await fetch("http://localhost:8090/favorites/ethereum", {
  method: "DELETE",
  headers: {
    Authorization: "Bearer my-secret-token",
  },
});

const body = await response.json();
```

---

### HTTP Status Codes

| Code | Meaning                                                          |
| ---- | ---------------------------------------------------------------- |
| 200  | Request succeeded (GET, DELETE)                                  |
| 201  | Resource created (POST)                                          |
| 400  | Bad request — missing or invalid request body or path parameter  |
| 401  | Unauthorized — missing, malformed, or invalid Bearer token       |
| 404  | Not found — the specified favorite does not exist for this user  |
| 409  | Conflict — the cryptocurrency is already in the user's favorites |
| 500  | Internal server error — database query or scan failure           |

---

## Deployment Modes

### Standalone

Run `favorites-api` independently with its own PostgreSQL instance.

> **Requires cpp-rest-api running at `:8080`** — this service delegates JWT token validation to `cpp-rest-api`. Without it, all requests will fail authentication. Start `cpp-rest-api` first (either natively or via its own `docker compose up`), then run:

```bash
# From the favorites-api/ directory
docker compose up --build
```

- Favorites API available at `http://localhost:8090`
- PostgreSQL available at `localhost:5432`
- `AUTH_API_URL` is set to `http://host.docker.internal:8080` in the compose file, which resolves to the host machine's `cpp-rest-api` from inside Docker (works on Docker Desktop for Windows and macOS)

### Full Stack (via crypto-dashboard)

When deployed as part of the full stack, `favorites-api` is orchestrated by `crypto-dashboard`'s `docker-compose.yml`. In this mode it shares a PostgreSQL instance with `cpp-rest-api` and both run on the same Docker network, so `AUTH_API_URL` is set to `http://cpp-rest-api:8080` and resolves correctly via Docker DNS.

See the [crypto-dashboard repository](https://github.com/lgarciac1603/crypto-dashboard) for full setup instructions.

---

## Running the Application

### Option 1: Local (go run)

```bash
# Set environment variables (or rely on defaults)
export DB_HOST=localhost
export DB_PORT=8090
export DB_NAME=apidb
export DB_USER=apiuser_test
export DB_PASS=apipass_test
export APP_PORT=8090

go run .
```

The server starts on `http://localhost:8090`.

### Option 2: Docker Compose (recommended)

This option starts both the PostgreSQL database and the API in isolated containers. The `init.sql` script automatically creates the `user_favorites` table on the first run.

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up --build -d

# Stop all services
docker-compose down

# Stop and remove volumes (resets the database)
docker-compose down -v
```

The API will be available at `http://localhost:8090` once the health check for PostgreSQL passes.

### Option 3: Compiled Binary

```bash
# Build a static binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o favorites-api .

# Run the binary
./favorites-api
```

---

## Testing

### Testing Framework

This project uses [Testify](https://github.com/stretchr/testify) (v1.11.1), the de-facto standard testing library for enterprise Go projects. Testify provides:

- `assert` — fluent assertion helpers that produce informative failure messages.
- `suite` — xUnit-style test suites with `SetupTest` / `TeardownTest` lifecycle hooks.

Database interactions in handler tests are isolated using [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock), which intercepts SQL queries without requiring a live database connection.

All tests follow the **Table-Driven Test** pattern and are organized into suites that mirror the package structure.

---

### Test Files and Coverage

| Package      | Test File                    | Test Cases | Coverage            |
| ------------ | ---------------------------- | ---------- | ------------------- |
| `config`     | `config/config_test.go`      | 5          | 100.0%              |
| `database`   | `database/database_test.go`  | 4          | 92.3%               |
| `handlers`   | `handlers/favorites_test.go` | 18         | 96.9%               |
| `middleware` | `middleware/auth_test.go`    | 8          | 100.0%              |
| `models`     | `models/favorite_test.go`    | 5          | N/A (no statements) |

> Note: The `database` package tests require a reachable PostgreSQL instance with the credentials defined in the [Configuration](#configuration) section. Tests for the remaining packages run fully offline.

---

### Running Tests

**Run all tests**

```bash
go test ./...
```

**Run all tests with verbose output**

```bash
go test -v ./...
```

**Run tests with basic coverage summary**

```bash
go test -cover ./...
```

**Generate a detailed HTML coverage report**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This opens an interactive report in your default browser highlighting covered and uncovered lines per file.

**Run tests for a single package**

```bash
go test -v ./handlers/...
go test -v ./middleware/...
go test -v ./config/...
```

---

### Test Cases Covered

**config package**

- Loads configuration values from environment variables.
- Falls back to default values when environment variables are not set.
- Generates a correctly formatted PostgreSQL connection string.
- Handles passwords containing special characters.
- Validates that all struct fields are populated.

**database package**

- Successfully opens and pings a live PostgreSQL connection.
- Returns a descriptive error when the host is unreachable.
- Returns an error when credentials are invalid.
- Handles `CloseDB` gracefully when no connection has been established.

**handlers package**

- `GetFavorites`: returns the full favorites list for a user.
- `GetFavorites`: returns an empty array when the user has no favorites.
- `GetFavorites`: returns `401` when `userID` is not present in context.
- `GetFavorites`: returns `500` when the database query fails.
- `GetFavorites`: returns `500` when a row scan fails.
- `PostFavorite`: creates a new favorite and returns `201`.
- `PostFavorite`: returns `400` when the request body is missing required fields.
- `PostFavorite`: returns `401` when `userID` is not present in context.
- `PostFavorite`: returns `409` when the cryptocurrency is already a favorite.
- `PostFavorite`: returns `500` when the insert query fails.
- `DeleteFavorite`: removes a favorite and returns `200`.
- `DeleteFavorite`: returns `401` when `userID` is not present in context.
- `DeleteFavorite`: returns `404` when the cryptocurrency is not in the user's favorites.
- `DeleteFavorite`: returns `500` when the delete query fails.

**middleware package**

- Returns `401` when the `Authorization` header is absent.
- Returns `401` when the header format is invalid (not two space-separated parts).
- Returns `401` when the `Bearer` keyword is missing.
- Returns `401` when the token value is empty after the `Bearer` prefix.
- Returns `401` when `Bearer` is lowercase.
- Returns `401` when the header contains double spaces.
- Extracts the token string and sets it as `userID` in the Gin context.
- Calls `c.Next()` after successful token extraction.

**models package**

- Marshals a `Favorite` struct to the expected JSON key names.
- Unmarshals JSON into a `Favorite` struct with correct field types.
- Validates `json` struct tag mappings (camelCase keys, no snake_case keys).
- Confirms field types (`int` for `ID` and `UserID`, `string` for others).
- Handles marshaling of a zero-value `Favorite` struct.

---

### Optional: Makefile

If you prefer short commands, add a `Makefile` to the project root:

```makefile
.PHONY: test test-verbose coverage coverage-html

test:
	go test ./...

test-verbose:
	go test -v ./...

coverage:
	go test -cover ./...

coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
```

Usage:

```bash
make test
make coverage
make coverage-html
```

---

## Integration with cpp-rest-api

Favorites API is designed to operate alongside **cpp-rest-api**, a C++ REST service that owns user management and authentication. The two services share the same PostgreSQL instance but maintain strict table-level separation.

### Architecture

```
+------------------------------------------+
|      cpp-rest-api (C++) -- PRIMARY        |
|  Responsibilities:                        |
|  - User registration and login            |
|  - JWT issuance                           |
|  - Token validation endpoint              |
|  - Core business logic                    |
+------------------------------------------+
              |
              | (trusts tokens issued by)
              v
+------------------------------------------+
|   favorites-api (Go) -- AUXILIARY         |
|  Responsibilities:                        |
|  - Validate Bearer tokens via cpp-rest    |
|  - Manage user_favorites table            |
|  - Scoped GET, POST, DELETE operations    |
+------------------------------------------+
              |
              | (both connect to)
              v
+------------------------------------------+
|      PostgreSQL (Shared Database)         |
|  - Table: users       (cpp-rest-api)      |
|  - Table: user_favorites (favorites-api)  |
+------------------------------------------+
```

### Communication Flow

```
1. Frontend authenticates via cpp-rest-api:
   POST /api/login  -->  cpp-rest-api  -->  returns JWT

2. Frontend sends a request to favorites-api with the JWT:
   GET /favorites
   Authorization: Bearer <jwt>

3. favorites-api middleware validates the token:
   POST /api/validate-token  -->  cpp-rest-api
   (current implementation accepts all non-empty tokens as a placeholder)

4. If valid, favorites-api:
   - Extracts userID from token or context
   - Executes a scoped query on user_favorites
   - Returns the result
```

### Required Change in cpp-rest-api

The only recommended addition to cpp-rest-api is the `user_favorites` table creation, placed in the database migration sequence:

```sql
-- database/migrations/005-user-favorites.sql
BEGIN;

CREATE TABLE IF NOT EXISTS user_favorites (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    crypto_id  VARCHAR(100) NOT NULL,
    crypto_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, crypto_id)
);

CREATE INDEX IF NOT EXISTS idx_user_favorites_user_id ON user_favorites(user_id);

COMMIT;
```

No changes to the `users` table, authentication logic, or existing endpoints in cpp-rest-api are required.

---

## Database Schema

### Table: user_favorites

```sql
CREATE TABLE IF NOT EXISTS user_favorites (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER      NOT NULL,
    crypto_id   VARCHAR(100) NOT NULL,
    crypto_name VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, crypto_id)
);
```

| Column        | Type         | Constraints                         | Description                                |
| ------------- | ------------ | ----------------------------------- | ------------------------------------------ |
| `id`          | SERIAL       | PRIMARY KEY                         | Auto-incrementing row identifier           |
| `user_id`     | INTEGER      | NOT NULL                            | References the authenticated user          |
| `crypto_id`   | VARCHAR(100) | NOT NULL                            | Cryptocurrency identifier (e.g. `bitcoin`) |
| `crypto_name` | VARCHAR(255) | NOT NULL                            | Human-readable name (e.g. `Bitcoin`)       |
| `created_at`  | TIMESTAMP    | NOT NULL, DEFAULT CURRENT_TIMESTAMP | Row creation timestamp                     |

The `UNIQUE (user_id, crypto_id)` constraint prevents duplicate entries at the database level, complementing the application-level `409 Conflict` check in `PostFavorite`.

**Sample Data**

| id  | user_id | crypto_id | crypto_name | created_at          |
| --- | ------- | --------- | ----------- | ------------------- |
| 1   | 1       | bitcoin   | Bitcoin     | 2026-04-02 10:30:00 |
| 2   | 1       | ethereum  | Ethereum    | 2026-04-02 10:35:00 |
| 3   | 2       | solana    | Solana      | 2026-04-02 11:00:00 |

---

## Authentication Flow

```
Client                      favorites-api              (future) cpp-rest-api
  |                               |                              |
  |-- GET /favorites ------------->|                              |
  |   Authorization: Bearer <tok> |                              |
  |                               |                              |
  |                        [AuthMiddleware]                      |
  |                          1. Read header                      |
  |                          2. Split on " "                     |
  |                          3. Validate "Bearer" prefix         |
  |                          4. ValidateToken(token)             |
  |                               |--- POST /api/validate-token ->|
  |                               |<-- { userId: 1 } ------------|
  |                          5. c.Set("userID", userID)          |
  |                          6. c.Next()                         |
  |                               |                              |
  |                        [GetFavorites handler]                |
  |                          7. c.Get("userID")                  |
  |                          8. SELECT ... WHERE user_id = $1    |
  |<-- 200 { data: [...] } --------|                              |
```

**Current Middleware Behavior (Placeholder)**

The `ValidateToken` function in `middleware/auth.go` currently accepts any non-empty string as a valid token and returns the raw token string as the `userID`. This is an intentional simplification to allow development without a running cpp-rest-api instance.

**Production Requirement**

Before deploying to production, `ValidateToken` must be replaced with a real HTTP call to cpp-rest-api's token validation endpoint. The returned `userID` must be an integer that matches a row in the `users` table, and handlers must cast accordingly.

---

## Best Practices

- **Dependency injection**: `FavoritesHandler` receives `*sql.DB` through its constructor (`NewFavoritesHandler`), keeping handlers testable without a real database connection.
- **Interface-based mocking**: The `DatabaseInterface` in the `database` package defines the subset of `*sql.DB` methods used by the application, enabling clean mocking in tests.
- **Parameterized queries**: All SQL statements use positional placeholders (`$1`, `$2`) to prevent SQL injection.
- **Scoped queries**: Every database query is filtered by `user_id`, ensuring strict data isolation between users.
- **Duplicate prevention**: The application checks for existing favorites before insertion and returns `409 Conflict`, reducing unnecessary writes and providing clear client feedback.
- **Graceful empty responses**: `GetFavorites` returns `[]` (never `null`) when no favorites exist, so consumers do not need nil-checks.
- **Environment-based configuration**: No credentials or hostnames are hard-coded. All values are read from environment variables with safe defaults for development.
- **Multi-stage Docker build**: The `Dockerfile` uses a builder stage to compile the binary and a minimal `alpine` image to run it, keeping the final image small and free of build tooling.
- **Non-root container user**: The production Docker image runs under a dedicated `appuser` account, reducing the attack surface.

---

## Contributing

Contributions are welcome. Please follow the guidelines below to keep the codebase consistent.

### Workflow

1. Fork the repository and create a feature branch from `main`:
   ```bash
   git checkout -b feat/short-description
   ```
2. Make your changes with focused, atomic commits.
3. Run the full test suite and confirm all tests pass:
   ```bash
   go test ./...
   ```
4. Open a pull request against `main` with a clear description of the change and its motivation.

### Commit Message Convention

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <short summary>
```

| Type       | When to use                                             |
| ---------- | ------------------------------------------------------- |
| `feat`     | A new feature                                           |
| `fix`      | A bug fix                                               |
| `docs`     | Documentation changes only                              |
| `test`     | Adding or correcting tests                              |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `chore`    | Dependency updates, build configuration, tooling        |
| `perf`     | Performance improvement                                 |

Examples:

```
feat(handlers): add pagination support to GetFavorites
fix(middleware): reject tokens with lowercase Bearer prefix
docs(readme): add architecture diagram for cpp-rest-api integration
test(handlers): cover PostFavorite 500 path with sqlmock
```

### Code Style

- Follow standard Go formatting: run `gofmt -w .` before committing.
- Keep functions small and focused on a single responsibility.
- Add tests for every new handler, middleware, or utility function.
- Do not commit `.env` files or files containing secrets.

---

## License

This project is licensed under the MIT License.

---

## Contact

Maintained by **lgarciac1603**.
Repository: [https://github.com/lgarciac1603/favorites-api](https://github.com/lgarciac1603/favorites-api)

---

## Roadmap

| Item                                    | Status  |
| --------------------------------------- | ------- |
| Core CRUD endpoints (GET, POST, DELETE) | Done    |
| PostgreSQL integration                  | Done    |
| Testify-based test suite                | Done    |
| Docker Compose setup                    | Done    |
| Real JWT validation via cpp-rest-api    | Pending |
| Pagination for GET /favorites           | Pending |
| Rate limiting                           | Pending |
| OpenAPI / Swagger documentation         | Pending |
| GitHub Actions CI pipeline              | Pending |
| Kubernetes deployment manifests         | Pending |
