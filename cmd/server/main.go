package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql-go-handler"
	"github.com/inkah-trace/server"
	"golang.org/x/net/context"
)

func main() {
	// define GraphQL schema using relay library helpers
	schema, err := graphql.NewSchema(server.Schema)

	if err != nil {
		log.Fatal(err)
	}

	h := NewCORSHandler(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	fs := http.FileServer(http.Dir("static"))

	// serve HTTP
	http.Handle("/graphql", h)
	http.Handle("/", fs)
	fmt.Println("Server listening on 127.0.0.1:9820")
	err = http.ListenAndServe("127.0.0.1:9820", nil)
	if err != nil {
		log.Fatal(err)
	}

}

type CORSHandler struct {
	graphQLGoHandler *handler.Handler
}

func (c CORSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")

	c.graphQLGoHandler.ContextHandler(context.Background(), w, r)
}

func NewCORSHandler(p *handler.Config) *CORSHandler {
	if p == nil {
		p = handler.NewConfig()
	}
	if p.Schema == nil {
		panic("undefined GraphQL schema")
	}

	return &CORSHandler{
		graphQLGoHandler: &handler.Handler{
			Schema: p.Schema,
		},
	}
}
