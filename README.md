### Overview

Drone telemetry pipeline: HTTP ingest pushes GPS, state, and events to Redpanda; a consumer worker persists them into TimescaleDB. A separate API server exposes REST endpoints and (planned) WebSocket for dashboards and live maps. Run ingestion-gateway, consumer-worker, and api as separate processes; use the Python simulator or any HTTP client to send telemetry.

### Getting started

1. Start TimescaleDB:
   ```bash
   docker compose up -d timescaledb
   ```

2. Run migrations (requires [golang-migrate](https://github.com/golang-migrate/migrate#installation) CLI):
   ```bash
   migrate -path ./migrations -database "postgres://drone:drone@localhost:5432/drone_platform?sslmode=disable" up
   ```

3. Generate types (requires [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) CLI). Ensure `$HOME/go/bin` is on your PATH (or run by full path):
   ```bash
   sqlc generate
   ```
   Output is written to `internal/platform/db/sqlc/`. No database connection is required. Re-run whenever you change `internal/platform/db/schema/` or `internal/platform/db/queries/`.

**Design:** Time-series tables (gps, state, events) use `(drone_id, time)` as the logical key. No surrogate insert ID is used; append-only telemetry is queried by drone and time range.

### TODO

1. Visualizations
	1. Get data to the frontend without having the database act as a middleman
	2. Create a new websocket Go server, or add to the existing one, that maintains a websocket connection with the frontend
	3. Every time a message comes from kafka, the server should send it to the db and push it out to the websockets - _Learning Focus:_ Channels in Go (how different parts of your code talk to each other safely).
2. Productionization
	1. Use a local kubernetes tool like `kind` or `minikube`
	2. Write the kubernetes manifests to deploy cluster