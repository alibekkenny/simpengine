---
description: Create a new golang-migrate migration pair with the next sequence number
argument-hint: <snake_case_name>
---

Create a new migration pair under `migrations/` named `$ARGUMENTS`.

## Steps

1. Inspect `migrations/` and determine the next sequence number (current highest + 1, zero-padded to 6 digits — e.g. if the latest is `000010_*`, the new one is `000011`).

2. Run:
   ```bash
   migrate create -ext sql -dir migrations -seq $ARGUMENTS
   ```
   This will create both `.up.sql` and `.down.sql` files with the next sequence number.

3. Write the `up` SQL. Use the same style as the existing migrations: lowercase keywords or uppercase consistently, snake_case columns, explicit foreign keys, sensible defaults, `created_at TIMESTAMPTZ DEFAULT NOW()` where applicable.

4. Write a real `.down.sql` that fully reverses the up migration. Don't leave it empty — that breaks `migrate down`.

5. Apply against the running dockerized DB to verify:
   ```bash
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/simpengine?sslmode=disable" up
   ```

6. Verify the down works too:
   ```bash
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/simpengine?sslmode=disable" down 1
   ```
   Then re-apply `up` so the schema is current.

Migrations run automatically on container start via `entrypoint.sh`, so anyone pulling will pick this up next `docker-compose up`.
