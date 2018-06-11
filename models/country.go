package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/lib/pq"
)

type Country struct {
	Code   string
	Name   string
	Cities []int64
}

type CountryPage struct {
	Countries []Country
}

func (db DB) Country(ctx context.Context, id string) (Country, error) {
	countries, err := db.Countries(ctx, []string{id})
	if err != nil {
		return Country{}, err
	}
	if len(countries) != 1 {
		return Country{}, errors.New("got wrong number of results")
	}
	return countries[0], nil
}

func (db DB) Countries(ctx context.Context, ids []string) ([]Country, error) {
	DBConn, err := sql.Open("postgres", db.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	var countries []Country

	query_string := `
	SELECT country.code, country.name,
	ARRAY_AGG(city.id) as city_id
	FROM country
	JOIN city ON (country.code = city.country_code)
	WHERE code = ANY($1)
	GROUP BY country.code
	`

	start := time.Now()
	rows, err := DBConn.Query(query_string, pq.Array(ids))
	elapsed := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var name string
		cityIDs := pq.Int64Array{}
		if err := rows.Scan(&code, &name, &cityIDs); err != nil {
			log.Fatal(err)
		}
		countries = append(countries, Country{code, name, []int64(cityIDs)})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("FETCH countries took %s for code(s) %s", elapsed, ids)
	return countries, nil
}
