package mngsession

import (
	"errors"
	"github.com/labstack/echo"
	"goappuser/models"
)

var (
	ErrNoSessionFound        = errors.New("No Session found")
	ErrSessionNotSessionizer = errors.New("Session Type is not a sessionizer")
)

type SessionUserReader interface {
	ReadSessionUser(c echo.Context, s models.Sessionizer) error
}

type SessionUserWriter interface {
	WriteSessionUser(c echo.Context, user models.User) error
}

type SessionUserIo interface {
	SessionUserReader
	SessionUserWriter
}
