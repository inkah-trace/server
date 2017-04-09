package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/inkah-trace/server/bolt"
	"log"
	"net/http"
	"github.com/inkah-trace/server/api"
	"fmt"
	"github.com/rs/cors"
)

func main() {
	bdbc := bolt.NewClient()
	bdbc.Path = "./inkah.boltdb"
	bdbc.Open()
	defer bdbc.Close()

	router := NewRouter()
	boltdb_router := AddBoltDBContext(router, bdbc)
	handler := cors.Default().Handler(boltdb_router)

	fmt.Printf("Starting Inkah server on port %d...\n", 50052)
	log.Fatal(http.ListenAndServe(":50052", handler))
}

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	for _, route := range api.Routes {
		switch route.Method {
		case "GET":
			router.GET(route.Pattern, route.Handle)
		case "POST":
			router.POST(route.Pattern, route.Handle)
		}
	}

	return router
}

func AddBoltDBContext(next *httprouter.Router, client *bolt.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "boltdb", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}