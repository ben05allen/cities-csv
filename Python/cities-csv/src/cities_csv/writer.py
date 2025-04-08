from typing import TYPE_CHECKING, Any

from .models import City

if TYPE_CHECKING:
    from sqlalchemy.orm import Session


def write_cities(session: "Session", cities: list[dict[str, Any]]):
    session.bulk_insert_mappings(City, cities)
