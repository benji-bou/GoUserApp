package middlewares

import (
	"github.com/labstack/echo"
	"log"
	"net/http"
)

// "github.com/gorilla/context"

func loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(r.URL)
		next(c)
	}
}
