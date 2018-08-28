package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/users/{id}",
		GetUserRouteHandler,
	},
	Route{
		"POST",
		"/users",
		PostUserRouteHandler,
	},
	Route{
		"PATCH",
		"/users/{id}",
		PatchUserRouteHandler,
	},
	Route{
		"POST",
		"/users/{id}/increase-age",
		PostIncreaseUserAgeRouteHandler,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandlerFunc)

	}
	return router
}
