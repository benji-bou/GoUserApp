package middlewares

import (
	"fmt"
	"goappuser/user"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func NewAuthorizationMiddleware(roles ...string) MiddlewareHandler {
	return AuthorizationMiddlewareHandler{roles: roles}
}

type AuthorizationMiddlewareHandler struct {
	roles []string
}

func (a AuthorizationMiddlewareHandler) Middle(w http.ResponseWriter, r *http.Request, next func()) {
	session := context.Get(r, "Session").(*sessions.Session)
	if usr, isOk := session.Values["user"].(user.User); isOk {
		for _, elem := range a.roles {
			if elem == usr.Role {
				next()
				return
			}
		}
	}
	fmt.Fprint(w, http.StatusForbidden)
}
