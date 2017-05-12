package middlewares

import (
	"github.com/labstack/echo"
	"goappuser/manager/mngsession"
	"goappuser/models"
	"log"
	"reflect"
)

func MiddlewareWithSessionType(sessionType models.Sessionizer, sm mngsession.SessionReader) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			inter := reflect.New(reflect.TypeOf(sessionType).Elem()).Interface()
			log.Println("Middleware - Session", reflect.TypeOf(inter))
			if session, isOK := inter.(models.Sessionizer); isOK == true {
				err := sm.ReadSession(c, session)
				if err != nil && err != mngsession.ErrNoSessionFound {
					log.Println("session error Middleware", err)
				} else if err != mngsession.ErrNoSessionFound {
					log.Println("Session found")
					c.Set("Session", session)
				} else {
					log.Println(mngsession.ErrNoSessionFound)
				}
			} else {
				log.Println("session error Middleware --> Problem getting a sessionizer")
			}
			next(c)
			return nil
		}
	}
}
