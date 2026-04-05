# Drone Telemetry Platform

A real-time drone telemetry ingestion and query platform built in Go. The system accepts high-frequency GPS, state, and event data from drone agents over HTTP, streams it through a Kafka-compatible message bus, persists it in a time-series database, and exposes it via a REST/WebSocket API for live dashboards and historical queries.

---

## Architecture

```
┌─────────────────┐     HTTP      ┌───────────────────┐     Kafka      ┌──────────────────┐
│ Drone / Simulator│──────────────▶│ ingestion-gateway │───────────────▶│  consumer-worker │
└─────────────────┘               └───────────────────┘                └────────┬─────────┘
                                                                                 │ writes
                                                                                 ▼
                                  ┌───────────────────┐               ┌──────────────────────┐
                                  │    api-server      │◀──────────────│     TimescaleDB       │
                                  │  REST + WebSocket  │    queries    │  (gps, state, events) │
                                  └───────────────────┘               └──────────────────────┘
```

The platform is composed of three independent Go processes:

**`ingestion-gateway`** — HTTP server that receives telemetry payloads from drones or simulators and publishes them to the appropriate Redpanda topic (`v1_gps`, `v1_state`, `v1_events`). Stateless and horizontally scalable.

**`consumer-worker`** — Kafka consumer that reads from all three topics and persists records into TimescaleDB. Owns the write path; responsible for data integrity and ordering guarantees.

**`api-server`** — REST API (with planned WebSocket support) that serves queries over the persisted telemetry. Exposes endpoints for time-range queries, drone state lookups, and event history. WebSocket push for live dashboard feeds is in progress.

---

## Data Model

Telemetry is split into three append-only time-series tables, each using `(drone_id, time)` as the logical key. No surrogate insert IDs — records are identified by drone and timestamp, which is the correct pattern for time-series workloads and enables efficient range scans in TimescaleDB.

| Table | Fields |
|-------|--------|
| `gps` | `drone_id`, `time`, `lat`, `lon`, `altitude`, `heading`, `speed` |
| `state` | `drone_id`, `time`, `battery_pct`, `flight_mode`, `armed`, `signal_strength` |
| `events` | `drone_id`, `time`, `event_type`, `payload` |

TimescaleDB hypertables provide automatic time-based partitioning, enabling fast time-range queries without manual partition management.

Type-safe database access is generated via [sqlc](https://sqlc.dev) — queries live in `internal/platform/db/queries/` and the generated Go code in `internal/platform/db/sqlc/` is never edited by hand.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23 |
| Message Bus | Redpanda (Kafka-compatible) |
| Database | TimescaleDB (PostgreSQL 16) |
| DB Query Gen | sqlc |
| Migrations | golang-migrate |
| Kafka Client | segmentio/kafka-go |
| Postgres Driver | lib/pq |
| Simulator | Python 3 |
| Infra | Docker Compose |

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and Docker Compose
- [Go 1.23+](https://go.dev/dl/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate#installation)
- [sqlc CLI](https://docs.sqlc.dev/en/latest/overview/install.html)

### 1. Start infrastructure

Brings up TimescaleDB, Redpanda, and the Redpanda Console. The `redpanda-init` container automatically creates the `v1_gps`, `v1_state`, and `v1_events` topics on startup.

```bash
docker compose up -d
```

Health checks are configured on both TimescaleDB and Redpanda — dependent services wait until both are ready before starting.

### 2. Run migrations

```bash
migrate -path ./migrations \
  -database "postgres://drone:drone@localhost:5432/drone_platform?sslmode=disable" \
  up
```

### 3. Generate typed DB code

Only required after modifying schemas or queries in `internal/platform/db/`. Ensure `$HOME/go/bin` is on your `PATH`.

```bash
sqlc generate
```

Output is written to `internal/platform/db/sqlc/`. No database connection is required.

### 4. Run the services

Each service runs as a separate process. Open three terminal windows:

```bash
# Terminal 1 — ingest gateway (HTTP receiver)
go run ./cmd/ingestion-gateway

# Terminal 2 — consumer worker (Kafka → TimescaleDB)
go run ./cmd/consumer-worker

# Terminal 3 — API server
go run ./cmd/api
```

### 5. Send telemetry

Use the Python simulator to generate and stream mock drone telemetry to the ingest gateway:

```bash
python scripts/drone-simulator/simulator.py
```

Or send a manual HTTP request:

```bash
curl -X POST http://localhost:8081/v1/telemetry/gps \
  -H "Content-Type: application/json" \
  -d '{"drone_id": "drone-001", "lat": 37.7749, "lon": -122.4194, "altitude": 120.5, "heading": 270.0, "speed": 15.2}'
```

### 6. Inspect topics (optional)

Redpanda Console is available at [http://localhost:8080](http://localhost:8080) — use it to inspect topic throughput, consumer lag, and individual messages.
