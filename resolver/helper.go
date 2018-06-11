package resolver

import (
	graphql "github.com/graph-gophers/graphql-go"
)

func extractID(id string) graphql.ID {
	return graphql.ID(id)
}

func strValue(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func nullableStr(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}
