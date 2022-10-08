package graph

import (
	"FICSIT-Ordis/internal/domain/ordis"
	"FICSIT-Ordis/internal/ports/gql/graph/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	O *ordis.Ordis
}

var Directives generated.DirectiveRoot

type DirectiveRoot struct {
	generated.DirectiveRoot
}
