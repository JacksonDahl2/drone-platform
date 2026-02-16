"""Drone load simulator: mocks many drones posting to the ingestion gateway."""

import asyncio
import logging

import click
from factory import create_drones
from ids import generate_new_ids, read_ids_from_file
from runner import run

logging.basicConfig(level=logging.INFO)


@click.command()
@click.option(
    "--drones",
    prompt="Number of drones",
    default=5,
    help="Number of drones to simulate",
)
@click.option(
    "--gateway",
    prompt="gateway",
    default="http://localhost:3000",
    help="Url of the ingestion gateway",
)
@click.option("--duration", help="Duration of the run, leave blank to run indefinitely")
@click.option("--ids-file", help="PATH to the ids file")
def simulate_drones(
    drones: int = 5,
    gateway: str = "http://localhost:3000",
    duration: int | None = None,
    ids_file: str | None = None,
):
    """Program that will simulate drones producing data to a api server."""
    if duration is not None:
        try:
            duration_sec = int(duration)
        except ValueError:
            raise click.BadParameter("duration must be a number (seconds) or omitted")
    else:
        duration_sec = None

    try:
        ids = (
            read_ids_from_file(drones, ids_file)
            if ids_file
            else generate_new_ids(drones)
        )
    except (FileNotFoundError, PermissionError, OSError) as e:
        raise click.FileError(ids_file or "", str(e))
    except ValueError as e:
        raise click.BadParameter(str(e))

    drone_list = create_drones(ids)
    try:
        asyncio.run(run(drone_list, gateway, duration_sec))
    except KeyboardInterrupt:
        click.echo("Stopped.", err=True)
        raise SystemExit(130)
    except Exception as e:
        logging.exception("run failed")
        raise click.ClickException(str(e))


if __name__ == "__main__":
    simulate_drones()
