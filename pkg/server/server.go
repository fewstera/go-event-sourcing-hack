package server

import (
	"net/http"

	"github.com/fewstera/go-event-sourcing-hack/pkg/user"
	"github.com/sirupsen/logrus"
)

// Server is a struct for the server. A struct is used so that the main program
// can pass in dependencies, such as a logger, the projection and command handler.
type Server struct {
	projection *user.Projection
	cmdHandler *user.CommandHandler
	log        *logrus.Logger
}

// NewServer creates a new Server.
func NewServer(log *logrus.Logger, projection *user.Projection, cmdHandler *user.CommandHandler) *Server {
	s := new(Server)
	s.projection = projection
	s.cmdHandler = cmdHandler
	s.log = log
	return s
}

// Start makes the server begin serving requests.
func (s *Server) Start() error {
	return http.ListenAndServe(":8000", s.router())
}
