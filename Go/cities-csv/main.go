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

type Config struct {
	CSVPath string
	DBPath  string
}

type Record []string

type City struct {
	Name       string  `validate:"required"`
	State      string  `validate:"required"`
	Population *uint32 `validate:"omitempty"`
	Latitude   float64 `validate:"latitude,required"`
	Longitude  float64 `validate:"longitude,required"`
}

func parseCity(r Record) (*City, error) {
	var populationPtr *uint32
	if r[2] != "" {
		population_64, err := strconv.ParseUint(r[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid population '%s': %v", r[2], err)
		}
		population_32 := uint32(population_64)
		populationPtr = &population_32
	}

	latitude, err := strconv.ParseFloat(r[3], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid latitude '%s': %v", r[3], err)
	}

	longitude, err := strconv.ParseFloat(r[4], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid longitude '%s': %v", r[4], 64)
	}

	c := City{Name: r[0], State: r[1], Population: populationPtr, Latitude: latitude, Longitude: longitude}
	return &c, nil
}

func main() {
	config := Config{CSVPath: "../../data/cities.csv", DBPath: "../../data/cities.db"}

	file, err := os.Open(config.CSVPath)
	if err != nil {
		log.Fatal("Error opening csv file:", err)
		return
	}
	defer file.Close()

	db, err := sql.Open("sqlite3", config.DBPath)
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

			city, err := parseCity(record)
			if err != nil {
				log.Printf("Error parsing CSV record: %v", err)
			}

			if err := validate.Struct(city); err != nil {
				log.Printf("Validation failed for record %v: %v", record, err)
			}

			citiesChan <- *city
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
