package app

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	// "net/url"
	// "tools"
)

//Middleware function type for Middleware
// type Middleware func(w Response, r *Request, next func())

//AppLoader function to load server before launch it
type AppLoader func() http.Handler

//ServerConfig struct for configuring server
type ServerConfig struct {
	Host         string
	Port         string
	SSL          bool
	SessionStore sessions.Store
}

//Server struct of the server
type Server struct {
	config *ServerConfig
}

func newServer(config *ServerConfig) *Server {

	return &Server{config: config}
}

//Start start http server
func Start(config *ServerConfig, loader AppLoader) {
	s := newServer(config)
	handler := loader()
	if config.SSL == false {
		log.Fatal(http.ListenAndServe(s.config.Host+":"+s.config.Port, handler))
	} else {
		// log.Fatal(http.ListenAndServe(s.config.Host+":"+s.config.Port, s.routes.Router))
		log.Fatal(http.ListenAndServeTLS(s.config.Host+":"+s.config.Port, "cert/server.pem", "cert/server.key", handler))
	}
}
