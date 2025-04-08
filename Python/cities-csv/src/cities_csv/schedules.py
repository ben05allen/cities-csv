from pydantic import BaseModel, Field, ConfigDict, to_pascal


class City(BaseModel):
    name: str = Field(alias="City")
    state: str
    population: int | None
    latitude: float
    longitude: float

    model_config = ConfigDict(alias_generator=to_pascal)
