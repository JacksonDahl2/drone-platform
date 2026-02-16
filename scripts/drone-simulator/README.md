# Drone simulator

## How to run

From `scripts/drone-simulator`:

```bash
uv sync
uv run python main.py
```

With options:

```bash
uv run python main.py --drones 10 --gateway http://localhost:3000 --duration 60
uv run python main.py --drones 5 --ids-file drone_ids.csv
```

| Option       | Default               | Description |
|-------------|------------------------|-------------|
| `--drones`  | 5                      | Number of drones to simulate. |
| `--gateway` | http://localhost:3000  | Base URL of the ingestion gateway. |
| `--duration`| (none)                 | Run for this many seconds. Omit to run until Ctrl+C. |
| `--ids-file`| (none)                 | Path to file with drone IDs (CSV or one per line). If omitted, IDs are generated as UUIDs. |

Have the ingestion gateway running on the `--gateway` URL before starting.
