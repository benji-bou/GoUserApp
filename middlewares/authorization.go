package middlewares

import (
	"errors"
	"github.com/labstack/echo"
	sessions "goappuser/session"
	"goappuser/user"
	"log"
)

func NewAuthorizationMiddleware(roles ...string) *AuthorizationMiddlewareHandler {
	return &AuthorizationMiddlewareHandler{roles: roles}
}

type AuthorizationMiddlewareHandler struct {
	roles []string
}

func (a *AuthorizationMiddlewareHandler) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("auth with roles accepted", a.roles)
		session := c.Get("Session").(*sessions.Session)
		log.Println("session found in context ", session)
		usr := session.User.(*user.User)
		log.Println("user session role", usr.Role)
		for _, elem := range a.roles {
			if elem == usr.Role {
				log.Println("auth ok", usr.Email)
				next(c)
				return nil
			}
		}
		return errors.New("User not authorize")
	}
}
