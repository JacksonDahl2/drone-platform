from drone import Drone


def create_drone(id: str) -> Drone:
    """Builds a single Drone with deterministic state from its id."""
    return Drone(id)


def create_drones(ids: list[str]) -> list[Drone]:
    """Builds one Drone per id; same id list yields same drone set."""
    return [Drone(drone_id) for drone_id in ids]
