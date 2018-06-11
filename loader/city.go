package loader

import (
	"context"
	"log"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/aywrite/gql/errors"
	"github.com/aywrite/gql/models"
)

func LoadCity(ctx context.Context, id string) (models.City, error) {
	var city models.City

	ldr, err := extract(ctx, cityLoaderKey)
	if err != nil {
		return city, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(id))()
	if err != nil {
		return city, err
	}

	city, ok := data.(models.City)
	if !ok {
		return city, errors.WrongType(city, data)
	}

	return city, nil
}

func LoadCities(ctx context.Context, ids []string) (CityResults, error) {
	var results []CityResult
	ldr, err := extract(ctx, cityLoaderKey)
	if err != nil {
		return results, err
	}

	data, errs := ldr.LoadMany(ctx, dataloader.NewKeysFromStrings(ids))()
	for i, d := range data {
		var e error
		if errs != nil {
			e = errs[i]
		}

		city, ok := d.(models.City)
		if !ok && e == nil {
			e = errors.WrongType(city, d)
		}

		results = append(results, CityResult{City: city, Error: e})
	}

	return results, nil
}

// CityResult is the (data, error) pair result of loading a specific key.
type CityResult struct {
	City  models.City
	Error error
}

// CityResults is a named type, so methods can be attached to []CityResult.
type CityResults []CityResult

// WithoutErrors filters any result pairs with non-nil errors.
func (results CityResults) WithoutErrors() []models.City {
	var cities = make([]models.City, 0, len(results))

	for _, r := range results {
		if r.Error != nil {
			continue
		}

		cities = append(cities, r.City)
	}

	return cities
}

func PrimeCities(ctx context.Context, page models.CityPage) error {
	ldr, err := extract(ctx, cityLoaderKey)
	if err != nil {
		return err
	}

	for _, c := range page.Cities {
		ldr.Prime(ctx, dataloader.StringKey(c.ID), c)
	}

	return nil
}

type cityGetter interface {
	City(ctx context.Context, id int) (models.City, error)
	Cities(ctx context.Context, id []int) ([]models.City, error)
}

// FilmLoader contains the client required to load film resources.
type cityLoader struct {
	get cityGetter
}

func newCityLoader(client cityGetter) dataloader.BatchFunc {
	return cityLoader{get: client}.loadBatch
}

func (ldr cityLoader) loadBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var ids []int
	for _, k := range keys {
		val, err := strconv.Atoi(k.String())
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, val)
	}

	cities, err := ldr.get.Cities(ctx, ids)
	if err != nil {
		log.Fatal(err)
		// return a `[]*dataloader.Result` with `len(ids)` elements that have errors.
	}

	// do any work you need to do to get the response in the same order as the input `ids`

	results := []*dataloader.Result{}
	for _, c := range cities {
		results = append(results, &dataloader.Result{c, nil})
	}
	return results
}
