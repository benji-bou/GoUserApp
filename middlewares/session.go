package middlewares

import (

	// "github.com/gorilla/context"
	"fmt"
	"goappuser"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

var sessionStore sessions.Store

//NewSessionMiddleware create a Session Middleware
func NewSessionMiddleware(store sessions.Store) MiddlewareHandlerFunc {
	sessionStore = store
	return sessionMiddleware
}

func sessionMiddleware(w http.ResponseWriter, r *http.Request, next func()) {
	session, err := sessionStore.Get(r, "session-key")
	log.Println("in session")
	if err != nil {
		log.Println("Middleware Session error = ", err)
		app.JSONResp(w, app.RequestError{"Session", "Cannot retrieve session", 0})
		fmt.Fprint(w, http.StatusInternalServerError)
		log.Println("in session return")
		return
	}
	log.Println(session)
	context.Set(r, "Session", session)
	next()
}
