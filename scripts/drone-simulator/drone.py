import hashlib
import math
import random
from faker import Faker

_METERS_PER_DEG_LAT = 111320.0


class Drone:
    """Holds mutable telemetry state for one simulated drone. Deterministic per id via seeded random."""

    def __init__(self, drone_id: str) -> None:
        """Seeds state from drone_id so the same id yields reproducible telemetry."""
        hash_seed = hashlib.sha256(drone_id.encode()).hexdigest()
        seed_int = int(hash_seed[:16], 16)
        self.id = drone_id
        self._r = random.Random(seed_int)

        Faker.seed(seed_int)
        fake = Faker()

        self.lat = float(fake.latitude())
        self.lon = float(fake.longitude())
        self.alt = self._r.uniform(0.0, 150.0)
        self.heading = self._r.uniform(0.0, 360.0)
        self.speed = self._r.uniform(0.0, 15.0)
        self.pitch = self._r.uniform(-5.0, 5.0)
        self.roll = self._r.uniform(-5.0, 5.0)
        self.climb_rate = self._r.uniform(-1.0, 1.0)
        self.angular_rate = self._r.uniform(-2.0, 2.0)
        self.battery_pct = self._r.uniform(70.0, 100.0)
        self.voltage = 15.0
        self.status = self._r.choice(("idle", "flying", "armed", "landing"))
        self.connected = True
        self.flight_mode = self._r.choice(("manual", "auto", "loiter"))
        self.waypoint_index = 0
        self.mission_id = str(fake.uuid4()) if self._r.random() < 0.6 else None

    def next_gps(self) -> None:
        """Advances position and attitude by one tick; mutates lat, lon, alt, heading, etc."""
        dt = 1.0
        rad = math.radians(self.heading)
        meters_lat = self.speed * math.cos(rad) * dt
        meters_lon = self.speed * math.sin(rad) * dt
        self.lat += meters_lat / _METERS_PER_DEG_LAT
        self.lon += meters_lon / (
            _METERS_PER_DEG_LAT * max(0.01, math.cos(math.radians(self.lat)))
        )
        self.alt = max(0.0, self.alt + self.climb_rate * dt)
        self.heading = (self.heading + self.angular_rate * dt) % 360.0
        self.speed = max(0.0, min(25.0, self.speed + self._r.uniform(-0.5, 0.5)))
        self.pitch = max(-15.0, min(15.0, self.pitch + self._r.uniform(-1.0, 1.0)))
        self.roll = max(-15.0, min(15.0, self.roll + self._r.uniform(-1.0, 1.0)))
        self.climb_rate = max(
            -3.0, min(3.0, self.climb_rate + self._r.uniform(-0.2, 0.2))
        )
        self.angular_rate = max(
            -5.0, min(5.0, self.angular_rate + self._r.uniform(-0.3, 0.3))
        )
        self.lat = max(-90.0, min(90.0, self.lat))
        self.lon = (self.lon + 180.0) % 360.0 - 180.0

    def next_state(self) -> None:
        """Advances battery, voltage, status, connected, flight_mode; may flip randomly."""
        self.battery_pct = max(0.0, self.battery_pct - self._r.uniform(0.02, 0.08))
        self.voltage = 15.0 * (self.battery_pct / 100.0)
        if self._r.random() < 0.03:
            self.status = self._r.choice(("idle", "flying", "armed", "landing"))
        if self._r.random() < 0.02:
            self.connected = not self.connected
        if self._r.random() < 0.02:
            self.flight_mode = self._r.choice(("manual", "auto", "loiter"))

    def maybe_event(self) -> tuple[str, dict] | None:
        """Returns (event_type, payload) when an event fires, else None."""
        if self.battery_pct < 20 and self._r.random() < 0.3:
            return ("low_battery", {"battery_pct": self.battery_pct})
        if self.battery_pct < 5 and self._r.random() < 0.5:
            return ("critical_battery", {"battery_pct": self.battery_pct})
        if self.status == "landing" and self._r.random() < 0.2:
            return ("mission_ended", {})
        if self.status == "flying" and self._r.random() < 0.08:
            event = (
                "waypoint_reached",
                {
                    "waypoint_index": self.waypoint_index,
                    "lat": self.lat,
                    "lon": self.lon,
                },
            )
            self.waypoint_index += 1
            return event
        if self.status == "flying" and self.mission_id and self._r.random() < 0.05:
            return ("mission_started", {"mission_id": self.mission_id})
        if not self.connected and self._r.random() < 0.25:
            return ("connection_lost", {})
        return None
