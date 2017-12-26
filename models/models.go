package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
	graphql "github.com/neelance/graphql-go"
)

type City struct {
	ID          int
	Name        string
	CountryCode string
}

type Country struct {
	Code   string
	Name   string
	Cities []int64
}

type Language struct {
	ID   string
	Name string
}

func RetrieveCities(ctx context.Context, ids []graphql.ID) ([]*City, error) {
	var cities []*City

	connStr := "user=world dbname=world-db password=world123 sslmode=disable"
	DBConn, err := sql.Open("postgres", connStr)
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
		cities = append(cities, &City{id, name, code})
	}
	if err := rows.Err(); err != nil {
		return cities, err
	}

	log.Printf("FETCH cites took %s for code(s) = %s", elapsed, ids)
	return cities, nil
}

func RetrieveCountries(ctx context.Context, ids []graphql.ID) ([]*Country, error) {
	connStr := "user=world dbname=world-db password=world123 sslmode=disable"
	DBConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	var countries []*Country

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
		countries = append(countries, &Country{code, name, []int64(cityIDs)})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("FETCH countries took %s for code(s) %s", elapsed, ids)
	return countries, nil
}
