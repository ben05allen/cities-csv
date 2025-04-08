package main() 

import (
  "database/sql"
  "encoding/csv"
  "fmt"
  "os"

  _ "github.com/mattn/go-sqlite3"
)

func main() {
  file, err := os.Open("cities.csv")
  if err != nil {
    fmt.Println("Error opening csv file:", err)
    return
  }
  defer file.Close()

  reader := csv.NewReader(file)
  records, err := reader.ReadAll()
  if err != nil {
    fmt.Println("Error reading csv file:", err)
    return
  }

  db, err := sql.Open("sqlite3", data.db)
  if err != nil {
    fmt.Println("Error opening database: ", err)
    return
  }
  defer db.Close()

  _, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS cities (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    state NOT NULL,
    population,
    latitude NOT NULL,
    longituge NOT NULL
  )
  `)
  if err != nil {
    fmt.Println("Error creating table:", err)
    return
  }


}

