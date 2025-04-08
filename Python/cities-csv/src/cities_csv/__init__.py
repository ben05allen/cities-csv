from csv import DictReader
from pathlib import Path

from .schedules import City


def main() -> None:
    csv_file = Path(__file__).parents[4] / "data" / "cities.csv"

    with open(csv_file) as f:
        dict_reader = DictReader(f)

        for row in dict_reader:
            print(City(**row).model_dump())
