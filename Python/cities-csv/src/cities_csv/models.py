# pyright: basic


from sqlalchemy import Integer
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column


class Base(DeclarativeBase):
    pass


class City(Base):
    __tablename__ = "cities"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    name: Mapped[str]
    state: Mapped[str]
    population: Mapped[int | None]
    latitude: Mapped[float]
    longitude: Mapped[float]
