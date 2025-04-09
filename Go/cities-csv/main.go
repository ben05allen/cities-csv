package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
)

type City struct {
	Name       string  `validate:"required"`
	State      string  `validate:"required"`
	Population *uint32 `validate:"omitempty"`
	Latitude   float64 `validate:"latitude,required"`
	Longitude  float64 `validate:"longitude,required"`
}

func main() {
	file, err := os.Open("../../data/cities.csv")
	if err != nil {
		log.Fatal("Error opening csv file:", err)
		return
	}
	defer file.Close()

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
    longitude REAL NOT NULL
  )
  `)
	if err != nil {
		log.Fatal("Error creating table:", err)
		return
	}

	citiesChan := make(chan City, 100)
	var wg sync.WaitGroup
	validate := validator.New()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := csv.NewReader(file)
		// skip header row
		_, err := reader.Read()
		if err != nil {
			log.Fatal("Error reading CSV headers:", err)
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading CSV record: %v", err)
				continue
			}

			var populationPtr *uint32
			if record[2] != "" {
				population_64, err := strconv.ParseUint(record[2], 10, 32)
				if err != nil {
					log.Printf("Invalid population '%s': %v", record[2], err)
					continue
				}
				population_32 := uint32(population_64)
				populationPtr = &population_32
			}

			latitude, err := strconv.ParseFloat(record[3], 64)
			if err != nil {
				log.Printf("Invalid latitude '%s': %v", record[3], err)
			}

			longitude, err := strconv.ParseFloat(record[4], 64)
			if err != nil {
				log.Printf("Invalid longitude '%s': %v", record[4], 64)
			}

			city := City{Name: record[0], State: record[1], Population: populationPtr, Latitude: latitude, Longitude: longitude}
			if err := validate.Struct(city); err != nil {
				log.Printf("Validation failed for record %v: %v", record, err)
			}

			citiesChan <- city
		}
		close(citiesChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		insertQuery := `INSERT INTO cities (name, state, population, latitude, longitude) VALUES (?, ?, ?, ?, ?)`
		stmt, err := db.Prepare(insertQuery)
		if err != nil {
			log.Fatal("Error preparing insert statement:", err)
		}
		defer stmt.Close()

		for city := range citiesChan {
			_, err := stmt.Exec(city.Name, city.State, city.Population, city.Latitude, city.Longitude)
			if err != nil {
				log.Printf("Error inserting record: %v", err)
				continue
			}
		}
	}()

	// wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Successfully imported CSV to SQLite!")

}
