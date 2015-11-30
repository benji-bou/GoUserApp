package middlewares

import (
	"goappuser/user"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

func NewAuthorizationMiddleware(roles ...string) MiddlewareHandler {
	return &AuthorizationMiddlewareHandler{roles: roles}
}

type AuthorizationMiddlewareHandler struct {
	roles []string
}

func (a *AuthorizationMiddlewareHandler) Middle(w http.ResponseWriter, r *http.Request, next func()) {
	log.Println("auth with roles accepted", a.roles)
	session := context.Get(r, "Session").(*sessions.Session)
	if usr, isOk := session.Values["user"].(user.User); isOk {
		log.Println("user session role", usr.Role)
		for _, elem := range a.roles {
			if elem == usr.Role {
				next()
				return
			}
		}
	}
	http.Error(w, "Vous n'êtes pas authentifié", http.StatusForbidden)
}
