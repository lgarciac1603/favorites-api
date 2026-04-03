# Favorites API

Simple REST API built with Go, Gin, and PostgreSQL to manage a user's favorite cryptocurrencies.

## Overview

This service exposes protected endpoints under `/favorites`:

- List favorites for the authenticated user
- Add a cryptocurrency to favorites
- Remove a cryptocurrency from favorites

The project is organized by responsibility:

- `main.go`: app bootstrap, route wiring, middleware, server start
- `config/`: environment-based configuration
- `database/`: PostgreSQL initialization and lifecycle
- `handlers/`: HTTP handlers for favorites
- `middleware/`: authentication middleware
- `models/`: API/domain models

## Tech Stack

- Go 1.26+
- Gin (`github.com/gin-gonic/gin`)
- PostgreSQL driver (`github.com/lib/pq`)
- Testing with `testing`, `testify`, and `sqlmock`

## API Endpoints

All routes are protected by `AuthMiddleware()` and require an `Authorization` header with format:

`Authorization: Bearer <token>`

### 1) Get Favorites

- Method: `GET`
- Path: `/favorites`

Successful response (`200`):

```json
{
  "data": [
    {
      "id": 1,
      "userId": 1,
      "cryptoId": "bitcoin",
      "cryptoName": "Bitcoin",
      "createdAt": "2026-04-02T10:30:00Z"
    }
  ],
  "total": 1
}
```

### 2) Add Favorite

- Method: `POST`
- Path: `/favorites`
- JSON body:

```json
{
  "cryptoId": "ethereum",
  "cryptoName": "Ethereum"
}
```

Possible responses:

- `201` created
- `400` missing/invalid payload
- `401` unauthenticated
- `409` already exists
- `500` database error

### 3) Delete Favorite

- Method: `DELETE`
- Path: `/favorites/:cryptoId`

Possible responses:

- `200` deleted
- `400` missing `cryptoId`
- `401` unauthenticated
- `404` not found
- `500` database error

## Authentication Behavior

Current middleware behavior is a placeholder:

- It validates header format (`Bearer <token>`)
- It treats any non-empty token as valid
- It stores the token value as `userID` in context

This is intentionally simple and should be replaced by a real integration with a users/auth service.

## Prerequisites

- Go installed
- PostgreSQL running and accessible

## Environment Variables

The app reads these variables (defaults shown):

- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_NAME=apidb`
- `DB_USER=apiuser_test`
- `DB_PASS=apipass_test`
- `APP_PORT=8090`

## Run Locally

1. Set environment variables.
2. Ensure PostgreSQL contains the required `user_favorites` table.
3. Start the API:

```bash
go run .
```

By default, the server starts on `:8090`.

## Test Suite

The repository includes unit/integration-style tests in all core packages:

- `config/config_test.go`
- `database/database_test.go`
- `handlers/favorites_test.go`
- `middleware/auth_test.go`
- `models/favorite_test.go`

### Run All Tests (Verbose)

```bash
go test -v ./...
```

Observed result (latest run):

- All tests passed
- Verbose output included all suites and subtests
- Root module has no test files

### Run Coverage

```bash
go test -cover ./...
```

Observed coverage (latest run):

- `github.com/lgarciac1603/favorites-api`: `0.0%` (no test files)
- `github.com/lgarciac1603/favorites-api/config`: `100.0%`
- `github.com/lgarciac1603/favorites-api/database`: `92.3%`
- `github.com/lgarciac1603/favorites-api/handlers`: `96.9%`
- `github.com/lgarciac1603/favorites-api/middleware`: `100.0%`
- `github.com/lgarciac1603/favorites-api/models`: `no statements`

## Notes

- `database` tests require a reachable PostgreSQL instance with matching credentials.
- Current auth middleware sets `userID` as a string token, while handlers cast `userID` to `int`. Tests pass because handler tests inject integer `userID` values directly. This mismatch should be aligned before production use.
