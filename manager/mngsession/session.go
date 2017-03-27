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

type SessionReader interface {
	ReadSession(c echo.Context, s models.Sessionizer) error
}

type SessionWriter interface {
	WriteSession(c echo.Context, s models.Sessionizer) error
}

type SessionUserIo interface {
	SessionReader
	SessionWriter
}
