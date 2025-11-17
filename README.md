# RMS API — quick reference

This document lists the API endpoints implemented so far, expected request/response shapes and short curl examples.

General notes
- The server reads configuration from `config/config.yml` (sections `config` and `auth`).
- On first run (database creation) an admin account is created using `config.admin_email` and `admin_password_hash`/`admin_password`.
- Authentication: Bearer JWT in the `Authorization: Bearer <TOKEN>` header (do not wrap the token in single quotes).
- Access tokens include the `role` claim (string) and `sub` (user id). Supported roles: `admin`, `editor`, `operator`, `mechanic`, `viewer`, `client`.

Run locally
```bash
cd rms
go run ./cmd/api
# or build and run
go build -o bin/rms-server ./cmd/api && ./bin/rms-server
```

Public endpoints

- GET /healthz
  - Returns DB status: `{ "status": "ok" }` or an error status.

- POST /auth/login
  - Body (JSON):
    ```json
    { "email": "user@example.com", "password": "secret" }
    ```
  - Response (200):
    ```json
    { "access_token": "...", "refresh_token": "..." }
    ```

- POST /auth/refresh
  - Body (JSON):
    ```json
    { "refresh_token": "<REFRESH_TOKEN>" }
    ```
  - Response (200): `{ "access_token": "..." }` — the new access token contains `role`.

- GET /workshop
  - Returns workshop info from config `car_workshop`:
    ```json
    { "name": "My Workshop", "description": "Auto repair shop", "image_path": "", "owner": "" }
    ```

Protected endpoints (require Bearer access token)

Authentication header: `Authorization: Bearer <ACCESS_TOKEN>`

1) Cars management — allowed roles: `admin`, `editor`, `operator`

- GET /cars
  - Returns array of `Car` objects.

- POST /cars
  - Body (JSON) example (fields map to DB columns):
    ```json
    {
      "owner_id": 1,
      "plate_number": "ABC1234",
      "vin": "1HGCM82633A004352",
      "make": "Toyota",
      "model": "Corolla",
      "year": 2015,
      "last_mileage": 120000,
      "fuel_type": "petrol",
      "engine_capacity": 1.8,
      "engine_type": "inline-4",
      "default_hourly_rate": 100.0,
      "notes": "Test car"
    }
    ```
  - Response: created `Car` object (201)

- PATCH /cars/:id
  - Body: partial `Car` fields to update. Response: updated object.

- DELETE /cars/:id
  - Deletes the car (204 No Content).

Example cURL to create a car:
```bash
ACCESS="<ACCESS_TOKEN>"
curl -X POST http://127.0.0.1:10800/cars \
  -H "Authorization: Bearer $ACCESS" \
  -H "Content-Type: application/json" \
  -d '{"owner_id":1, "plate_number":"ABC1234", "make":"Toyota", "model":"Corolla"}'
```

User account endpoints (authenticated)

- PATCH /me
  - Update own profile (e.g. `user_name`, `email`). Example body:
    ```json
    { "user_name": "NewName", "email": "me@example.com" }
    ```
  - Response: updated `User` JSON.

- POST /me/password
  - Change password. Body:
    ```json
    { "old_password": "old", "new_password": "new" }
    ```
  - Response: 204 No Content on success.

Tokens and claims
- Access tokens include `sub` (user id as string), `role` (string), `exp`, `iat`, `jti`.
- Refresh tokens are minimal `RegisteredClaims` (`sub`, `exp`, `iat`, `jti`).

Configuration (`config/config.yml`)
- `config.admin_email`, `config.admin_password` / `config.admin_password_hash` — initial admin account.
- `auth.jwt_key` — secret used to sign JWTs (keep it secure).
- `auth.jwt_ttl`, `auth.refresh_token_ttl` — token TTLs (e.g. "24h", "168h").
- `config.db_path` — path to sqlite DB file.

Debug tips
- If you get `invalid token`, ensure the Authorization header contains the raw token without surrounding quotes.
- If you get `forbidden`, verify the access token contains the `role` claim and the DB user role matches (`admin|editor|operator`).

Next steps I can help with:
- add detailed curl examples for every endpoint,
- generate an OpenAPI spec or Postman collection from the routes.
