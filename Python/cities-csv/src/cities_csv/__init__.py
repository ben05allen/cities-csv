# pyright: basic

import aiofiles
import asyncio
import aiosqlite
from csv import DictReader
from pathlib import Path

from .schedules import City


async def create_table(db):
    await db.execute("""
            CREATE TABLE IF NOT EXISTS cities (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                city TEXT NOT NULL,
                state TEXT NOT NULL,
                population INTEGER,
                latitude REAL NOT NULL,
                longitude REAL NOT NULL
            );
        """)
    await db.commit()


async def read_csv(file_path: Path):
    async with aiofiles.open(str(file_path.resolve())) as f:
        content = await f.read()

    lines = content.splitlines()

    reader = DictReader(lines)
    for row in reader:
        yield row


async def async_main() -> None:
    csv_file = Path(__file__).parents[4] / "data" / "cities.csv"
    db_file = Path(__file__).parents[4] / "data" / "cities.db"

    if db_file.exists():
        db_file.unlink()

    async with aiosqlite.connect(str(db_file.resolve())) as db:
        await create_table(db)

        async for row in read_csv(csv_file):
            parsed_row = City(**row).model_dump()  # pyright: ignore
            await db.execute(
                "INSERT INTO cities (city, state, population, latitude, longitude) VALUES (?, ?, ?, ?, ?)",
                tuple(parsed_row.values()),
            )

        await db.commit()


def main():
    asyncio.run(async_main())
