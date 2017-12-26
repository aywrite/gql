package world

import (
	"context"
	"errors"
	"gql/models"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	graphql "github.com/neelance/graphql-go"
	"github.com/nicksrandall/dataloader"
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

type Resolver struct{}

type CountryResolver struct {
	code   graphql.ID
	loader *dataloader.Loader
}

func NewCountryResolver(ctx context.Context, id graphql.ID) (*CountryResolver, error) {
	loader, found := ctx.Value("countryLoader").(*dataloader.Loader)
	if !found {
		return nil, errors.New("unable to find counrty loader")
	}

	if id == graphql.ID("") {
		return nil, errors.New("no county ID specified")
	}

	return &CountryResolver{id, loader}, nil
}

func (r *CountryResolver) load(ctx context.Context) (*models.Country, error) {
	// we can have any kinds of necessary checks here
	if r.loader == nil {
		return nil, errors.New("missing country loader")
	}

	// kind of verbose, but makes code bulletproof and easy to debug
	if r.code == graphql.ID("") {
		return nil, errors.New("missing country key")
	}

	// use the loader we attached in the constructor
	thunk := r.loader.Load(ctx, r.code)
	data, err := thunk()
	if err != nil {
		return nil, err
	}

	country, ok := data.(*models.Country)
	if !ok {
		return nil, errors.New("unable to convert response to Country")
	}
	return country, nil
}

func (r *CountryResolver) ID(ctx context.Context) graphql.ID {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return graphql.ID(c.Code)
}

func (r *CountryResolver) Name(ctx context.Context) string {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c.Name
}

func (r *CountryResolver) Cities(ctx context.Context) *[]*CityResolver {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var resolvers []*CityResolver
	for _, cityID := range c.Cities {
		id := graphql.ID(strconv.FormatInt(cityID, 10))
		resolver, err := NewCityResolver(ctx, id)
		if err != nil {
			log.Fatal(err)
		}
		resolvers = append(resolvers, resolver)
	}
	return &resolvers
}

func (r *Resolver) Country(ctx context.Context, args struct{ ID graphql.ID }) *CountryResolver {
	resolver, err := NewCountryResolver(ctx, args.ID)
	if err != nil {
		log.Fatal(err)
	}
	return resolver
}

type CityResolver struct {
	id     graphql.ID
	loader *dataloader.Loader
}

func NewCityResolver(ctx context.Context, id graphql.ID) (*CityResolver, error) {
	loader, found := ctx.Value("cityLoader").(*dataloader.Loader)
	if !found {
		return nil, errors.New("unable to find city loader")
	}

	if id == graphql.ID("") {
		return nil, errors.New("no city ID specified")
	}

	return &CityResolver{id, loader}, nil
}

func (r *CityResolver) load(ctx context.Context) (*models.City, error) {
	// we can have any kinds of necessary checks here
	if r.loader == nil {
		return nil, errors.New("missing city loader")
	}

	// kind of verbose, but makes code bulletproof and easy to debug
	if r.id == graphql.ID("") {
		return nil, errors.New("missing city key")
	}

	// use the loader we attached in the constructor
	thunk := r.loader.Load(ctx, r.id)
	data, err := thunk()
	if err != nil {
		return nil, err
	}

	city, ok := data.(*models.City)
	if !ok {
		return nil, errors.New("unable to convert response to City")
	}
	return city, nil
}

func (r *CityResolver) ID(ctx context.Context) graphql.ID {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return graphql.ID(strconv.Itoa(c.ID))
}

func (r *CityResolver) Name(ctx context.Context) string {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return c.Name
}

func (r *CityResolver) Country(ctx context.Context) *CountryResolver {
	c, err := r.load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	resolver, err := NewCountryResolver(ctx, graphql.ID(c.CountryCode))
	if err != nil {
		log.Fatal(err)
	}
	return resolver
}

func (r *Resolver) City(ctx context.Context, args struct{ ID graphql.ID }) *CityResolver {
	resolver, err := NewCityResolver(ctx, args.ID)
	if err != nil {
		log.Fatal(err)
	}
	return resolver
}
