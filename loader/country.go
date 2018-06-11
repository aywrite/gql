package loader

import (
	"context"
	"log"

	"github.com/graph-gophers/dataloader"

	"github.com/aywrite/gql/errors"
	"github.com/aywrite/gql/models"
)

func LoadCountry(ctx context.Context, code string) (models.Country, error) {
	var country models.Country

	ldr, err := extract(ctx, countryLoaderKey)
	if err != nil {
		return country, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(code))()
	if err != nil {
		return country, err
	}

	country, ok := data.(models.Country)
	if !ok {
		return country, errors.WrongType(country, data)
	}

	return country, nil
}

func LoadCountries(ctx context.Context, codes []string) (CountryResults, error) {
	var results []CountryResult
	ldr, err := extract(ctx, countryLoaderKey)
	if err != nil {
		return results, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(codes))()
	for i, d := range data {
		var e error
		if errs != nil {
			e = errs[i]
		}

		country, ok := d.(models.Country)
		if !ok && e == nil {
			e = errors.WrongType(country, d)
		}

		results = append(results, CountryResult{Country: country, Error: e})
	}

	return results, nil
}

// CountryResult is the (data, error) pair result of loading a specific key.
type CountryResult struct {
	Country models.Country
	Error   error
}

// CountryResults is a named type, so methods can be attached to []CountryResult.
type CountryResults []CountryResult

// WithoutErrors filters any result pairs with non-nil errors.
func (results CountryResults) WithoutErrors() []models.Country {
	var countries = make([]models.Country, 0, len(results))

	for _, r := range results {
		if r.Error != nil {
			continue
		}

		countries = append(countries, r.Country)
	}

	return countries
}

func PrimeCountries(ctx context.Context, page models.CountryPage) error {
	ldr, err := extract(ctx, countryLoaderKey)
	if err != nil {
		return err
	}

	for _, c := range page.Countries {
		ldr.Prime(ctx, dataloader.StringKey(c.Code), c)
	}

	return nil
}

type countryGetter interface {
	Country(ctx context.Context, code string) (models.Country, error)
	Countries(ctx context.Context, codes []string) ([]models.Country, error)
}

// countryLoader contains the client required to load country resources.
type countryLoader struct {
	get countryGetter
}

func newCountryLoader(client countryGetter) dataloader.BatchFunc {
	return countryLoader{get: client}.loadBatch
}

func (ldr countryLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var ids []string
	for _, k := range keys {
		ids = append(ids, k.String())
	}

	countries, err := ldr.get.Countries(ctx, ids)
	if err != nil {
		log.Fatal(err)
		// return a `[]*dataloader.Result` with `len(ids)` elements that have errors.
	}

	// do any work you need to do to get the response in the same order as the input `ids`

	results := []*dataloader.Result{}
	for _, c := range countries {
		results = append(results, &dataloader.Result{c, nil})
	}
	return results
}
