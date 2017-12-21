package world

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	graphql "github.com/neelance/graphql-go"
)

var BasicSchema = `
	schema {
		query: Query
	}
	type Query {
		city(id: ID!): City
		country(id: ID!): Country
    }
	type City {
		id: ID!
		name: String!
		country: Country
    }
	type Country {
		id: ID!
		name: String!
		cities: [City]
    }
`

var (
	DBConn *sql.DB
)

type BasicResolver struct{}

func (r *BasicResolver) City(args struct{ ID graphql.ID }) *cityResolver {
	city, err := getCity(args.ID)
	if err != nil {
		log.Fatal(err)
	}
	return &cityResolver{city}
}

func (r *BasicResolver) Country(args struct{ ID graphql.ID }) *countryResolver {
	country, err := getCountry(args.ID)
	if err != nil {
		log.Fatal(err)
	}
	return &countryResolver{country}
}

type cityResolver struct {
	c *city
}

func (r *cityResolver) ID() graphql.ID {
	return graphql.ID(strconv.Itoa(r.c.ID))
}

func (r *cityResolver) Name() string {
	return r.c.Name
}

func (r *cityResolver) Country() *countryResolver {
	country, err := getCountry(graphql.ID(r.c.CountryCode))
	if err != nil {
		log.Fatal(err)
	}
	return &countryResolver{country}
}

type countryResolver struct {
	c *country
}

func (r *countryResolver) ID() graphql.ID {
	return graphql.ID(r.c.Code)
}

func (r *countryResolver) Name() string {
	return r.c.Name
}

func (r *countryResolver) Cities() *[]*cityResolver {
	code := graphql.ID(r.c.Code)
	return resolveCities(code)
}

type city struct {
	ID          int
	Name        string
	CountryCode string
}

type country struct {
	Code string
	Name string
}

type language struct {
	ID   string
	Name string
}

func getCity(ID graphql.ID) (*city, error) {
	var err error
	connStr := "user=world dbname=world-db password=world123 sslmode=disable"
	DBConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var id int
	var name string
	var code string
	err = DBConn.QueryRow("SELECT id, name, country_code FROM city WHERE id = $1", ID).Scan(&id, &name, &code)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("FETCH city id = %d", id)
	return &city{id, name, code}, nil
}

func resolveCities(ID graphql.ID) *[]*cityResolver {
	connStr := "user=world dbname=world-db password=world123 sslmode=disable"
	DBConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	var cities []*cityResolver

	rows, err := DBConn.Query("SELECT id, name, country_code FROM city WHERE country_code = $1", ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			log.Fatal(err)
		}
		cities = append(cities, &cityResolver{&city{id, name, code}})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("FETCH cites code = %s", ID)
	return &cities
}

func getCountry(ID graphql.ID) (*country, error) {
	var err error
	connStr := "user=world dbname=world-db password=world123 sslmode=disable"
	DBConn, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	var code string
	var name string
	err = DBConn.QueryRow("SELECT code, name FROM country WHERE code = $1", ID).Scan(&code, &name)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("FETCH country code = %s", code)
	return &country{code, name}, nil
}
