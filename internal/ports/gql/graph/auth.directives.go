package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func init() {
	Directives.IsAuthenticated = IsAuthenticated
}
func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	//opCtx := graphql.GetOperationContext(ctx)
	//auth := opCtx.Headers.Get("Authorization")
	return next(ctx)
}
