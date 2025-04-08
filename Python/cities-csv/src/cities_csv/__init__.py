from csv import DictReader
from pathlib import Path

from sqlalchemy.orm import scoped_session, sessionmaker

from .db import get_session
from .parser import parse_city
from .writer import write_cities


def main() -> None:
    csv_file = Path(__file__).parents[4] / "data" / "cities.csv"

    with open(csv_file) as f:
        dict_reader = DictReader(f)

        cities = [parse_city(row) for row in dict_reader]

    SessionFactory = sessionmaker()
    ScopedSession = scoped_session(SessionFactory)

    with Session() as session:
        write_cities(session, cities)
