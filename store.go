package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB
var q *sql.Stmt

// func init() {
// 	var err error

// 	db, err = sql.Open("postgres", "postgres://media:DexLx3WM14uhbwrvk9cC@localhost/media?sslmode=disable")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	q, err = db.Prepare("INSERT INTO download (title, description, guid, size, data) values ($1, $2, $3, $4, $5)")

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func Store(i *Item) error {
	b, err := json.Marshal(i)

	if err != nil {
		return fmt.Errorf("Error marshaling item: %v", err)
	}

	size, _ := strconv.Atoi(i.Attrs["size"])

	//zero size should be considered okay
	_, err = q.Exec(i.Title, i.Description, i.Attrs["guid"], size, b)

	if err != nil {
		return fmt.Errorf("SQL query exec error: %v", err)
	}

	return nil
}
