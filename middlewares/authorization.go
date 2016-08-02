package middlewares

import (
	"errors"
	"github.com/labstack/echo"
	"goappuser/models"
	"log"
	"net/http"
)

var (
	ErrUserNotAuthorized = errors.New("User not authorize")
	ErrUserRolesNotMatch = errors.New("User doesn't have corrects roles")
	ErrNoSessionUser     = errors.New("No user session found")
)

func NewAuthorizationMiddleware(roles ...string) *AuthorizationMiddlewareHandler {
	return &AuthorizationMiddlewareHandler{roles: roles}
}

type AuthorizationMiddlewareHandler struct {
	roles []string
}

func (a *AuthorizationMiddlewareHandler) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// log.Println("auth with roles accepted", a.roles)
		session, isOk := c.Get("Session").(models.Session)
		if !isOk {
			log.Println("In Authorization Session not found")
			c.JSON(http.StatusUnauthorized, models.RequestError{Title: "Authorization Error", Description: ErrUserNotAuthorized.Error(), Code: 0})
			return ErrUserNotAuthorized
		}
		// log.Println("session found in context ", session)
		usr := session.User
		// log.Println("user session role", usr.GetRole())
		for _, elem := range a.roles {
			if elem == usr.GetRole() {
				// log.Println("auth ok", usr.GetUniqueLogin())
				next(c)
				return nil
			}
		}
		c.JSON(http.StatusUnauthorized, models.RequestError{Title: "Authorization Error", Description: ErrUserRolesNotMatch.Error(), Code: 0})
		return ErrUserRolesNotMatch
	}
}
