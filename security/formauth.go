package security

import (
	"errors"
	"github.com/labstack/echo"

	"golang.org/x/crypto/bcrypt"
)

//BasicAuthentication manage authentication with username and passwords
type FormAuthentication struct {
}

//GetCredentials log user
func (a FormAuthentication) GetCredentials(c echo.Context) (string, string, error) {
	login := c.FormValue("login")
	password := c.FormValue("password")
	if login == "" || password == "" {
		return "", "", errors.New("empty authentication fields")
	}
	return login, password, nil
}

//Compare set of password
func (a FormAuthentication) Compare(clearPassword, hashedpassword []byte) bool {
	errcmp := bcrypt.CompareHashAndPassword(hashedpassword, clearPassword)
	return errcmp == nil
}

//Hash password in order
func (a FormAuthentication) Hash(clearpassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(clearpassword, 0)

}
