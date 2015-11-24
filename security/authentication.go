package security

import "net/http"

//AuthType type of authentication
type AuthType int

const (
	Basic AuthType = 0 << iota
)

//AuthenticationProcesser interfacefor authenticate credentials
type AuthenticationProcesser interface {
	//Authenticate authenticate the usser with credentials and return the username decoded and password hashed
	GetCredentials(r *http.Request) (string, string, error)
	Hash(clearpassword []byte) ([]byte, error)
	Compare(clearPassword, realmPassword []byte) bool
}

//NewAuth create auth from type
func NewAuth(auth AuthType) AuthenticationProcesser {
	switch auth {
	case Basic:
		return &BasicAuthentication{}
	default:
		return nil
	}
}
