package webserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	GetHandlers   map[string]http.HandlerFunc
	PostHandlers  map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {

	return &WebServer{
		Router:        chi.NewRouter(),
		GetHandlers:   make(map[string]http.HandlerFunc),
		PostHandlers:  make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}

}

func (s *WebServer) AddHandler(method string, path string, handler http.HandlerFunc) {
	switch strings.ToUpper(method) {
	case "POST":
		s.PostHandlers[path] = handler
	case "GET":
		s.GetHandlers[path] = handler
	default:
		panic(fmt.Errorf("invalid method %s ", method))
	}
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.GetHandlers {
		s.Router.Get(path, handler)
	}
	for path, handler := range s.PostHandlers {
		s.Router.Post(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
