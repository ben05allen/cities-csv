# pyright: basic
from typing import Any
from .schedules import City


def parse_city(city: dict[str, Any]):
    return City(**city).model_dump()
