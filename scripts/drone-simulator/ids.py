# two functions,
# 1. read in from file, check if its csv, read it if it is, otherwise just read it in
# 2. generate number of random ones based on the drone number input


import csv
import uuid
from pathlib import Path


def read_ids_from_file(drone_num: int, file_name: str) -> list[str]:
    """Reads up to drone_num ids from a file; treats as CSV if extension is .csv."""
    path = Path(file_name)
    ids: list[str] = []
    if path.suffix.lower() == ".csv":
        with path.open(newline="") as f:
            reader = csv.reader(f)
            next(reader, None)
            for row in reader:
                ids.extend(cell.strip() for cell in row if cell.strip())
    else:
        with path.open() as f:
            ids = [line.strip() for line in f if line.strip()]
    return ids[:drone_num]


def generate_new_ids(drone_num: int) -> list[str]:
    """Generates unique ids from the uuid4 library"""
    if drone_num < 0:
        raise ValueError("drone_num must be non-negative")
    return [str(uuid.uuid4()) for _ in range(drone_num)]
