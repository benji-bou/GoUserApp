package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

type validator func(r *http.Request) bool

var validate validator

//NewAccessControlMiddleware create a middleware accessControl validator
//validator callback to check the validation
func NewAccessControlMiddleware(val validator) MiddlewareHandlerFunc {
	validate = val
	return accessControl
}

func accessControl(w http.ResponseWriter, r *http.Request, next func()) {
	if validate(r) == true {
		if origin := r.Header.Get("Origin"); origin != "" {
			log.Println("access control ok", r.URL)
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, Cookie, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == "OPTIONS" {
				log.Println("Options Method accepted")
				fmt.Fprint(w, http.StatusAccepted)
				return
			}
		}
	}
	next()
}
