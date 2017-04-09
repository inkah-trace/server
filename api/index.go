package api

import (
	"net/http"
	"fmt"
	"html"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
