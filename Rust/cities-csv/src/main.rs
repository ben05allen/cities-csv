use csv::Reader;
use serde::Deserialize;
use sqlx::sqlite::SqlitePoolOptions;
use std::{error::Error, path::Path};
use tokio::sync::mpsc;

#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
struct City {
    #[serde(rename = "City")]
    name: String,
    latitude: f64,
    longitude: f64,
    #[serde(deserialize_with = "csv::invalid_option")]
    population: Option<u32>,
    state: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let csv_path = Path::new(env!("CARGO_MANIFEST_DIR")).join("cities.csv");
    let db_path = format!(
        "sqlite://{}",
        Path::new(env!("CARGO_MANIFEST_DIR"))
            .join("cities.db")
            .display()
    );

    let pool = SqlitePoolOptions::new()
        .max_connections(5)
        .connect(&db_path)
        .await?;

    sqlx::query(
        "CREATE TABLE IF NOT EXISTS cities (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                state TEXT NOT NULL,
                population INTEGER,
                latitude REAL NOT NULL,
                longitude REAL NOT NULL
            );",
    )
    .execute(&pool)
    .await?;

    let mut rdr = Reader::from_path(&csv_path)?;
    let (tx, mut rx) = mpsc::channel::<City>(100);

    let db_handle = tokio::spawn({
        let pool = pool.clone();
        async move {
            while let Some(city) = rx.recv().await {
                sqlx::query(
                    "INSERT INTO cities (name, state, population, latitude, longitude) 
                    VALUES (?, ?, ?, ?, ?);",
                )
                .bind(&city.name)
                .bind(&city.state)
                .bind(city.population)
                .bind(city.latitude)
                .bind(city.longitude)
                .execute(&pool)
                .await
                .unwrap();
            }
        }
    });

    for result in rdr.deserialize() {
        let city: City = result?;
        tx.send(city).await?;
    }

    drop(tx);
    db_handle.await?;

    pool.close().await;
    println!("Data written to database");
    Ok(())
}
