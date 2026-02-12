This is the python script that will spin up all the drone instances to load test this kafka instance


httpx -> api requests, sync and async
click -> cli commands
faker -> fake drone ids and other stuff
tenacity -> retry with backoff 

## Plan
1. CLI (Click)
Options: --drones N (default e.g. 100), --gateway-url URL, --duration SECONDS (how long to run, or “forever”), --ids-file PATH (optional path to file with one ID per line).
Flow: Parse args → load or generate IDs → build drone set → run simulation.

2. Drone IDs
From file: If --ids-file is set, read lines (strip, skip empty), use as drone IDs. Expect N lines if you want exactly N drones, or use first N / repeat as needed.
Generated: If no file, generate N UUIDs at startup (e.g. uuid.uuid4() or uuid.uuid4().hex). Store in a list so the same run uses the same IDs.
Single source of truth: One list of strings (length = number of drones). Factory and sim both use this list.

3. Drone “factory”
Role: Given a list of IDs, produce “drone” objects (or dicts) that know their ID and can produce the next telemetry payloads.
Output: Something you can iterate: e.g. Drone(id, initial_state) with methods like next_gps(), next_state(), maybe_event().
No HTTP here: Factory only creates drones and their state; a separate “runner” or “orchestrator” calls the gateway.

4. Drone state and “normal” change
Per-drone state: Keep a small state object per drone, e.g.:
lat, lon, alt, heading, speed (for GPS)
battery_pct, status, connected (for state)
Optional: mission_id, waypoint_index for events.
Normal change:
GPS: Each tick, add small deltas (e.g. lat += speed * cos(heading) * dt in crude form, or add a small random drift). Clamp to sane bounds. This gives smooth movement instead of random jumps.
Battery: Decrease slowly over time (e.g. battery_pct -= 0.01 * dt), clamp at 0.
State: Keep status/connected stable most of the time; occasionally flip (e.g. “flying” → “landing” after a while, or connected false for a few ticks).
Events: From time to time (e.g. random or every K ticks), emit “waypoint_reached”, “low_battery” when below threshold, “mission_ended”, etc., based on current state.
Determinism (optional): Seed random per drone (e.g. random.Random(id) or hash of id) so the same ID list gives reproducible behavior.

5. file layout
```
scripts/drone-simulator/
  main.py           # Click CLI, load IDs, start runner
  ids.py            # load_ids_from_file() / generate_ids(n)
  drone.py          # Drone state + next_gps(), next_state(), maybe_event()
  factory.py        # create_drones(id_list) -> list[Drone]
  payloads.py       # build_gps_payload(drone), build_state_payload(drone), build_event_payload(drone, event_type)
  runner.py         # run N drones: loop or async, call gateway (httpx), use payloads from each drone
  pyproject.toml
```
factory.py: Takes the ID list, returns a list of Drone instances (each initialized with starting state).
drone.py: Holds mutable state and “normal” update logic; exposes next_gps(), next_state(), maybe_event() (and optionally tick(dt) to advance time).
payloads.py: Pure functions: given a drone (and maybe event type), return the dict that matches your Go API (same shape as GpsInput, StateInput, EventInput).

6. Runner
Async: One httpx.AsyncClient, one task per drone. Each task loop: get drone.next_gps() → POST to /gps, sleep (e.g. 1 s), repeat; same for state (e.g. every 5 s) and events when maybe_event() returns something.
Option B – Threads: Same idea with threading and httpx.Client (or requests).

Can also use drone ids file as a persistent data store maybe? between restarts