package resolver

import (
	"context"
	"log"
	"strconv"

	"github.com/aywrite/gql/errors"
	"github.com/aywrite/gql/loader"
	"github.com/aywrite/gql/models"
	graphql "github.com/graph-gophers/graphql-go"
)

type CityResolver struct {
	city models.City
}

type CityResolverArgs struct {
	City models.City
	ID   string
}

type CitiesResolverArgs struct {
	Page models.CityPage
	IDs  []string
}

func NewCityResolver(ctx context.Context, args CityResolverArgs) (*CityResolver, error) {
	var city models.City
	var err error

	switch {
	case args.City.Name != "":
		city = args.City
	case args.ID != "":
		city, err = loader.LoadCity(ctx, args.ID)
	default:
		err = errors.UnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &CityResolver{city: city}, nil
}

func NewCitiesResolver(ctx context.Context, args CitiesResolverArgs) (*[]*CityResolver, error) {
	err := loader.PrimeCities(ctx, args.Page)
	if err != nil {
		return nil, err
	}

	results, err := loader.LoadCities(ctx, append(args.IDs, args.Page.IDs()...))
	if err != nil {
		return nil, err
	}

	var (
		cities    = results.WithoutErrors()
		resolvers = make([]*CityResolver, 0, len(cities))
		errs      errors.Errors
	)

	for i, city := range cities {
		resolver, err := NewCityResolver(ctx, CityResolverArgs{City: city})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers = append(resolvers, resolver)
	}

	return &resolvers, errs.Err()

}

func (r *CityResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(strconv.Itoa(r.city.ID))
}

func (r *CityResolver) Name(ctx context.Context) string {
	return r.city.Name
}

func (r *CityResolver) Country(ctx context.Context) (*CountryResolver, error) {
	return NewCountryResolver(ctx, CountryResolverArgs{ID: r.city.CountryCode})
}

func (r *Resolver) City(ctx context.Context, args struct{ ID int32 }) *CityResolver {
	resolver, err := NewCityResolver(ctx, CityResolverArgs{ID: strconv.Itoa(int(args.ID))})
	if err != nil {
		log.Fatal(err)
	}
	return resolver
}
