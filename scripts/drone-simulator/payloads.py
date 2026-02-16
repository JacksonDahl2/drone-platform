import datetime
from drone import Drone


class Payload:
    """Builds JSON payloads for the ingestion gateway from drone state."""

    def _tick_timestamp(self) -> str:
        return (
            datetime.datetime.now(datetime.timezone.utc)
            .isoformat()
            .replace("+00:00", "Z")
        )

    def build_gps_payload(self, drone: Drone) -> dict[str, any]:
        """Returns a dict matching GpsInput for POST /gps."""
        return {
            "drone_id": drone.id,
            "timestamp": self._tick_timestamp(),
            "latitude": drone.lat,
            "longitude": drone.lon,
            "altitude": drone.alt,
            "heading": drone.heading,
            "pitch": drone.pitch,
            "roll": drone.roll,
            "speed": drone.speed,
            "climb_rate": drone.climb_rate,
            "angular_rate": drone.angular_rate,
        }

    def build_state_payload(self, drone: Drone) -> dict[str, any]:
        """Returns a dict matching StateInput for POST /state."""
        return {
            "drone_id": drone.id,
            "timestamp": self._tick_timestamp(),
            "status": drone.status,
            "battery_pct": drone.battery_pct,
            "voltage": drone.voltage,
            "connected": drone.connected,
            "flight_mode": drone.flight_mode,
        }

    def build_event_payload(
        self, drone: Drone, event: tuple[str, dict]
    ) -> dict[str, any]:
        """Returns a dict matching EventInput for POST /events."""
        event_type, payload = event
        return {
            "drone_id": drone.id,
            "timestamp": self._tick_timestamp(),
            "event_type": event_type,
            "payload": payload,
        }
