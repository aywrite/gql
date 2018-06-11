package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/lib/pq"
)

type City struct {
	ID          int
	Name        string
	CountryCode string
}

type CityPage struct {
	Cities []City
}

func (cp *CityPage) IDs() []string {
	var ids []string
	for _, c := range cp.Cities {
		ids = append(ids, strconv.Itoa(c.ID))
	}
	return ids
}

func (db DB) City(ctx context.Context, id int) (City, error) {
	cities, err := db.Cities(ctx, []int{id})
	if err != nil {
		return City{}, err
	}
	if len(cities) != 1 {
		return City{}, errors.New("got wrong number of results")
	}
	return cities[0], nil
}

func (db DB) Cities(ctx context.Context, ids []int) ([]City, error) {
	var cities []City
	DBConn, err := sql.Open("postgres", db.ConnStr)
	if err != nil {
		return cities, err
	}

	start := time.Now()
	rows, err := DBConn.Query("SELECT id, name, country_code FROM city WHERE id = ANY($1)", pq.Array(ids))
	elapsed := time.Since(start)
	if err != nil {
		return cities, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			log.Fatal(err)
		}
		cities = append(cities, City{id, name, code})
	}
	if err := rows.Err(); err != nil {
		return cities, err
	}

	log.Printf("FETCH cites took %s for code(s) = %d", elapsed, ids)
	return cities, nil
}
