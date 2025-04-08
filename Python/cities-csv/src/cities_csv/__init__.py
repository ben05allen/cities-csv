from csv import DictReader
from pathlib import Path

from .schedules import City
from .models import City as Model


def main() -> None:
    csv_file = Path(__file__).parents[4] / "data" / "cities.csv"

    with open(csv_file) as f:
        dict_reader = DictReader(f)

        for row in dict_reader:
            validated_row = City(**row).model_dump()
            model_row = Model(**validated_row)
