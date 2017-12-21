package world

import (
	"context"
	"gql/models"
	"log"

	graphql "github.com/neelance/graphql-go"
	"github.com/nicksrandall/dataloader"
)

// CountryLoader probably holds whatever we will use to load Things.
type CountryLoader struct{}

func (l *CountryLoader) Attach(ctx context.Context) dataloader.BatchFunc {
	return func(ctx context.Context, keys []interface{}) []*dataloader.Result {
		var ids []graphql.ID
		for _, i := range keys {
			ids = append(ids, i.(graphql.ID))
		}

		countries, err := models.RetrieveCountries(ctx, ids)
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
}

// CountryLoader probably holds whatever we will use to load Things.
type CityLoader struct{}

func (l *CityLoader) Attach(ctx context.Context) dataloader.BatchFunc {
	return func(ctx context.Context, keys []interface{}) []*dataloader.Result {
		var ids []graphql.ID
		for _, i := range keys {
			ids = append(ids, i.(graphql.ID))
		}

		countries, err := models.RetrieveCities(ctx, ids)
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
}
