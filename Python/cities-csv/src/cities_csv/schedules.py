from pydantic import BaseModel, Field, ConfigDict, field_validator
from pydantic.alias_generators import to_pascal


class City(BaseModel):
    name: str = Field(alias="City")
    state: str
    population: int | None
    latitude: float
    longitude: float

    model_config = ConfigDict(alias_generator=to_pascal)

    @field_validator("population", mode="before")
    def coerce_empty_to_none(value):
        if not value:
            return None

        return value
