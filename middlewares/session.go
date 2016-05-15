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
ruct {
	roles []string
}

func (a *AuthorizationMiddlewareHandler) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("auth with roles accepted", a.roles)
		session := context.Get(r, "Session").(*sessions.Session)
		log.Println("sessio