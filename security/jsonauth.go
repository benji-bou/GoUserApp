package security

import (
	"errors"
	"github.com/labstack/echo"
)

//BasicAuthentication manage authentication with username and passwords
type JsonAuthentication struct {
	FormAuthentication
}

//GetCredentials log user
func (a JsonAuthentication) GetCredentials(c echo.Context) (string, string, error) {
	credentials := struct {
		Login    string
		Password string
	}{}
	c.Bind(&credentials)
	if credentials.Login == "" || credentials.Password == "" {
		return "", "", errors.New("empty authentication fields")
	}
	return credentials.Login, credentials.Password, nil
}

//Compare set of password
