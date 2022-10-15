package gql

import (
	"FICSIT-Ordis/internal/domain/ordis"
	"FICSIT-Ordis/internal/ports/gql/graph"
	"FICSIT-Ordis/internal/ports/gql/graph/generated"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func Server(o *ordis.Ordis) error {
	port := os.Getenv("ORDIS_GQL_PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(o.Auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			O: o,
		},
		Directives: graph.Directives,
	}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return http.ListenAndServe(":"+port, router)
}
