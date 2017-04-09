package api

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Handle httprouter.Handle
}

type RouteList []Route

var Routes = RouteList{
	Route{"Index", "GET", "/", Index},
	Route{"EventCreate", "POST", "/event", EventCreate},
	Route{"TraceList", "GET", "/trace", TraceList},
	Route{"TraceGet", "GET", "/trace/:id", TraceGet},
}
