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

type CountryResolver struct {
	country models.Country
}

type CountryResolverArgs struct {
	Country models.Country
	ID      string
}
type CountriesResolverArgs struct {
	Page models.CountryPage
	IDs  []string
}

func NewCountryResolver(ctx context.Context, args CountryResolverArgs) (*CountryResolver, error) {
	var country models.Country
	var err error

	switch {
	case args.Country.Name != "":
		country = args.Country
	case args.ID != "":
		country, err = loader.LoadCountry(ctx, args.ID)
	default:
		err = errors.UnableToResolve
	}

	if err != nil {
		return nil, err
	}

	return &CountryResolver{country: country}, nil
}

func NewCountriesResolver(ctx context.Context, args CountriesResolverArgs) (*[]*CountryResolver, error) {
	results, err := loader.LoadCountries(ctx, args.IDs)
	if err != nil {
		return nil, err
	}

	var (
		countries = results.WithoutErrors()
		resolvers = make([]*CountryResolver, 0, len(countries))
		errs      errors.Errors
	)

	for i, country := range countries {
		resolver, err := NewCountryResolver(ctx, CountryResolverArgs{Country: country})
		if err != nil {
			errs = append(errs, errors.WithIndex(err, i))
		}

		resolvers = append(resolvers, resolver)
	}

	return &resolvers, errs.Err()

}

func (r *CountryResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(r.country.Code)
}

func (r *CountryResolver) Name(ctx context.Context) string {
	return r.country.Name
}

func (r *CountryResolver) Cities(ctx context.Context) (*[]*CityResolver, error) {
	list := make([]string, len(r.country.Cities))
	for i := range r.country.Cities {
		list[i] = strconv.Itoa(int(r.country.Cities[i]))
	}
	return NewCitiesResolver(ctx, CitiesResolverArgs{IDs: list})
}

func (r *Resolver) Country(ctx context.Context, args CountryResolverArgs) *CountryResolver {
	resolver, err := NewCountryResolver(ctx, args)
	if err != nil {
		log.Fatal(err)
	}
	return resolver
}
