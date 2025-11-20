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
  - Example:
    ```bash
    curl -i http://127.0.0.1:10800/healthz
    ```

- POST /auth/login
  - Body (JSON):
    ```json
    { "email": "user@example.com", "password": "secret" }
    ```
  - Response (200):
    ```json
    { "access_token": "...", "refresh_token": "..." }
    ```
    - Example:
      ```bash
      curl -X POST http://127.0.0.1:10800/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"admin@example.com","password":"secret"}'
      ```

- POST /auth/refresh
  - Body (JSON):
    ```json
    { "refresh_token": "<REFRESH_TOKEN>" }
    ```
  - Response (200): `{ "access_token": "..." }` — the new access token contains `role`.
    - Example:
      ```bash
      curl -X POST http://127.0.0.1:10800/auth/refresh \
        -H "Content-Type: application/json" \
        -d '{"refresh_token":"<REFRESH_TOKEN>"}'
      ```

- GET /workshop
  - Returns workshop info from config `car_workshop`:
    ```json
    { "name": "My Workshop", "description": "Auto repair shop", "image_path": "", "owner": "" }
    ```
  - Example:
    ```bash
    curl http://127.0.0.1:10800/workshop
    ```


Protected endpoints (require Bearer access token)

Authentication header: `Authorization: Bearer <ACCESS_TOKEN>`

Note: protected endpoints are mounted under role-specific prefixes in the running server. Common mounts are:

- `/operator/...` — routes for operator role (also used by editor/admin where applicable)
- `/admin/...` — admin routes
- user routes such as `/me` remain at the top level

1) Cars management — allowed roles: `admin`, `editor`, `operator`

- GET `/operator/cars` and `/admin/cars`
  - Returns array of `Car` objects.
  - Example (list cars):
    ```bash
    curl -H "Authorization: Bearer $ACCESS" http://127.0.0.1:10800/operator/cars
    ```

- POST `/operator/cars` and `/admin/cars`
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

    - Example (create car):
      ```bash
      ACCESS="<ACCESS_TOKEN>"
      curl -X POST http://127.0.0.1:10800/operator/cars \
        -H "Authorization: Bearer $ACCESS" \
        -H "Content-Type: application/json" \
        -d '{"owner_id":1,"plate_number":"ABC1234","vin":"1HGCM82633A004352","make":"Toyota","model":"Corolla","year":2015}'
      ```

- PATCH `/operator/cars/:id` (and `/admin/cars/:id`)
  - Body: partial `Car` fields to update. Response: updated object.
  - Example (patch):
    ```bash
    curl -X PATCH http://127.0.0.1:10800/operator/cars/12 \
      -H "Authorization: Bearer $ACCESS" \
      -H "Content-Type: application/json" \
      -d '{"notes":"Updated note"}'
    ```

- DELETE `/operator/cars/:id` (and `/admin/cars/:id`)
  - Deletes the car (204 No Content).
  - Example (delete):
    ```bash
    curl -X DELETE http://127.0.0.1:10800/operator/cars/12 \
      -H "Authorization: Bearer $ACCESS"
    ```

Example cURL to create a car (operator namespace):
```bash
ACCESS="<ACCESS_TOKEN>"
curl -X POST http://127.0.0.1:10800/operator/cars \
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
    - Example:
      ```bash
      curl -X PATCH http://127.0.0.1:10800/me \
        -H "Authorization: Bearer $ACCESS" \
        -H "Content-Type: application/json" \
        -d '{"user_name":"NewName","email":"me@example.com"}'
      ```

- POST /me/password
  - Change password. Body:
    ```json
    { "old_password": "old", "new_password": "new" }
    ```
  - Response: 204 No Content on success.
    - Example:
      ```bash
      curl -X POST http://127.0.0.1:10800/me/password \
        -H "Authorization: Bearer $ACCESS" \
        -H "Content-Type: application/json" \
        -d '{"old_password":"old","new_password":"new"}'
      ```

Tokens and claims
- Access tokens include `sub` (user id as string), `role` (string), `exp`, `iat`, `jti`.
- Refresh tokens are minimal `RegisteredClaims` (`sub`, `exp`, `iat`, `jti`).

Repair endpoints (roles: `admin`, `editor`, `operator`)

- POST /orders
  - Create a new repair order for a car.
  - Body example:
    ```json
    { "car_id": 123, "title": "Brake service", "description": "Replace pads", "total_cost": 200.0 }
    ```
  - Example curl:
    ```bash
    curl -X POST http://127.0.0.1:10800/orders \
      -H "Authorization: Bearer $ACCESS" \
      -H "Content-Type: application/json" \
      -d '{"car_id":123,"title":"Brake service","description":"Replace pads","total_cost":200}'
    ```

- POST /tasks
  - Create a new task inside a repair order (assign to mechanic).
  - Body example:
    ```json
    { "order_id": 456, "mechanic_id": 12, "title": "Replace brake pads", "hours": 1.5 }
    ```
  - Example curl:
    ```bash
    curl -X POST http://127.0.0.1:10800/tasks \
      -H "Authorization: Bearer $ACCESS" \
      -H "Content-Type: application/json" \
      -d '{"order_id":456,"mechanic_id":12,"title":"Replace brake pads","hours":1.5}'
    ```

Configuration (`config/config.yml`)
- `config.admin_email`, `config.admin_password` / `config.admin_password_hash` — initial admin account.
- `auth.jwt_key` — secret used to sign JWTs (keep it secure).
- `auth.jwt_ttl`, `auth.refresh_token_ttl` — token TTLs (e.g. "24h", "168h").
- `config.db_path` — path to sqlite DB file.

Debug tips
- If you get `invalid token`, ensure the Authorization header contains the raw token without surrounding quotes.
- If you get `forbidden`, verify the access token contains the `role` claim and the DB user role matches (`admin|editor|operator`).
