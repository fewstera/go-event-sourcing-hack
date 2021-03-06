package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (s *Server) router() *mux.Router {
	var routes = Routes{
		Route{
			"GET",
			"/users/{id}",
			s.getUserHandler,
		},
		Route{
			"GET",
			"/users",
			s.getAllUsersHandler,
		},
		Route{
			"POST",
			"/users",
			s.createUserHandler,
		},
		Route{
			"POST",
			"/users/{id}/deposit",
			s.depositHandler,
		},
		Route{
			"POST",
			"/users/{id}/withdraw",
			s.withdrawHandler,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(route.HandlerFunc)

	}
	return router
}
