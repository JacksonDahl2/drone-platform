import csv
import uuid

if __name__ == "__main__":
    with open("drone_ids.csv", "w", newline="") as csvfile:
        writer = csv.writer(csvfile, delimiter=",")
        writer.writerow(["id"])
        for _ in range(200):
            writer.writerow([str(uuid.uuid4())])
