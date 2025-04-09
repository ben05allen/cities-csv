package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file, err := os.Open("../../data/cities.csv")
	if err != nil {
		log.Fatal("Error opening csv file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Fatal("Error reading CSV headers:", err)
	}
	fmt.Println("Headers: ", headers)

	db, err := sql.Open("sqlite3", "../../data/cities.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS cities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    state TEXT NOT NULL,
    population INTEGER,
    latitude REAL NOT NULL,
    longituge REAL NOT NULL
  )
  `)
	if err != nil {
		log.Fatal("Error creating table:", err)
		return
	}

	insertQuery := `INSERT INTO cities (name, state, population, latitude, longituge) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		log.Fatal("Error preparing insert statement:", err)
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV row:", err)
		}

		_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4])
		if err != nil {
			log.Fatal("Error inserting record:", err)
		}
	}
	fmt.Println("Successfully import CSV to SQLite")
}
