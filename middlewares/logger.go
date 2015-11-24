package middlewares

import (
	"log"
	"net/http"
)

// "github.com/gorilla/context"

func loggerMiddleware(w http.ResponseWriter, r *http.Request, next func()) {
	log.Println(r.URL)
}
