package security

import (
	"encoding/base64"
	"errors"
	"github.com/labstack/echo"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//BasicAuthentication manage authentication with username and passwords
type BasicAuthentication struct {
}

//GetCredentials log user
func (a BasicAuthentication) GetCredentials(c echo.Context) (string, string, error) {
	s := strings.SplitN(c.Request().Header().Get("Authorization"), " ", 2)

	if len(s) != 2 || s[0] != "Basic" {
		return "", "", errors.New("Not Basic authentication challenge")
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return "", "", errors.New("Not base 64 encoding")
	}

	parts := strings.SplitN(string(b), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("Credentials malformed shall be username:password")
	}
	return parts[0], parts[1], nil
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
