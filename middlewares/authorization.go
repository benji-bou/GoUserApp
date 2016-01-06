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
	log.Println("session found in context ", session)
	if usr, isOk := session.Values["user"].(*user.User); isOk {
		log.Println("user session role", usr.Role)
		for _, elem := range a.roles {
			if elem == usr.Role {
				log.Println("auth ok", usr.Email)
				next()
				return
			}
		}
		log.Println("User not authorize", usr)
	} else {
		log.Println("User not found in session", usr)

	}
	http.Error(w, "Vous n'êtes pas authentifié", http.StatusForbidden)
}
