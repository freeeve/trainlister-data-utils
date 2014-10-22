package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "gopkg.in/cq.v1"
)

func main() {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	  MATCH (t:Train)-[sa:STOPS_AT]->(s:Stop)
	  return t.id, t.direction, t.type, coalesce(t.saturday, false), s.name, s.mile, sa.hour, sa.minute
   `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	f, err := os.Create("temp.csv")
	if err != nil {
		log.Fatal(err)
	}

	w := csv.NewWriter(f)
	for rows.Next() {
		var id int
		var direction string
		var trainType string
		var saturday bool
		var name string
		var mile float64
		var hour int
		var minute int
		err := rows.Scan(&id, &direction, &trainType, &saturday, &name, &mile, &hour, &minute)
		if err != nil {
			log.Fatal(err)
		}
		strs := []string{fmt.Sprintf("%d", id), direction, trainType, fmt.Sprintf("%v", saturday), name, fmt.Sprintf("%.2f", mile), fmt.Sprintf("%d", hour), fmt.Sprintf("%d", minute)}

		w.Write(strs)
	}
	w.Flush()
	f.Close()
}
