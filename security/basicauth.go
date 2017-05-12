package security

import (
	"errors"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

//BasicAuthentication manage authentication with username and passwords
type BasicAuthentication struct {
}

//GetCredentials log user
func (a BasicAuthentication) GetCredentials(c echo.Context) (string, string, error) {

	username, password, ok := c.Request().BasicAuth()
	if ok == false {
		return "", "", errors.New("Not Basic authentication challenge")
	}

	return username, password, nil
}

//Compare set of password
func (a BasicAuthentication) Compare(clearPassword, hashedpassword []byte) bool {
	errcmp := bcrypt.CompareHashAndPassword(hashedpassword, clearPassword)
	return errcmp == nil
}

//Hash password in order
func (a BasicAuthentication) Hash(clearpassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(clearpassword, 0)

}
