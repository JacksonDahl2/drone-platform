import asyncio
import logging
import time
from drone import Drone
from payloads import Payload
from tenacity import retry, stop_after_attempt, wait_exponential

import httpx

log = logging.getLogger(__name__)

_STATE_INTERVAL_TICKS = 5


async def _post(client: httpx.AsyncClient, url: str, json: dict) -> None:
    response = await client.post(url, json=json)
    response.raise_for_status()


@retry(
    stop=stop_after_attempt(3),
    wait=wait_exponential(multiplier=0.5, min=1, max=4),
    reraise=True,
)
async def _post_with_retry(client: httpx.AsyncClient, url: str, payload: dict) -> None:
    await _post(client, url, payload)


async def drone_worker(
    client: httpx.AsyncClient,
    drone: Drone,
    payload_builder: Payload,
    gateway: str,
    duration: int | None,
) -> None:
    start = time.monotonic()
    tick = 0
    while True:
        if duration is not None and (time.monotonic() - start) >= duration:
            break
        drone.next_gps()
        if tick % _STATE_INTERVAL_TICKS == 0:
            drone.next_state()
        event = drone.maybe_event()

        gps_payload = payload_builder.build_gps_payload(drone)
        gps_url = f"{gateway}/gps"
        try:
            await _post_with_retry(client, gps_url, gps_payload)
        except Exception as e:
            log.warning("gps post failed drone=%s: %s", drone.id, e)

        if tick % _STATE_INTERVAL_TICKS == 0:
            state_payload = payload_builder.build_state_payload(drone)
            state_url = f"{gateway}/state"
            try:
                await _post_with_retry(client, state_url, state_payload)
            except Exception as e:
                log.warning("state post failed drone=%s: %s", drone.id, e)

        if event:
            event_payload = payload_builder.build_event_payload(drone, event)
            event_url = f"{gateway}/events"
            try:
                await _post_with_retry(client, event_url, event_payload)
            except Exception as e:
                log.warning("events post failed drone=%s: %s", drone.id, e)

        tick += 1
        await asyncio.sleep(1)


async def run(drones: list[Drone], gateway: str, duration: int | None) -> None:
    payload_builder = Payload()
    async with httpx.AsyncClient() as client:
        tasks = [
            asyncio.create_task(
                drone_worker(client, d, payload_builder, gateway, duration)
            )
            for d in drones
        ]
        await asyncio.gather(*tasks)
