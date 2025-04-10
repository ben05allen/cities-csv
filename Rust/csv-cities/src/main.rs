use serde::Deserialize;
use std::{env, error::Error, ffi::OsString, process};

#[derive(Debug, Deserialize)]
#[serde(rename_all = "PascalCase")]
struct Record {
    #[serde(rename = "City")]
    name: String,
    latitude: f64,
    longitude: f64,
    population: Option<u64>,
    state: String,
}

fn main() {
    if let Err(err) = run() {
        println!("{}", err);
        process::exit(1);
    }
}

fn run() -> Result<(), Box<dyn Error>> {
    let file_path = get_first_arg()?;
    let mut rdr = csv::Reader::from_path(file_path)?;

    for result in rdr.deserialize() {
        let record: Record = result?;
        println!("{:?}", record);
    }
    Ok(())
}

fn get_first_arg() -> Result<OsString, Box<dyn Error>> {
    match env::args_os().nth(1) {
        None => Err(From::from("expected 1 argument, but got none")),
        Some(file_path) => Ok(file_path),
    }
}
