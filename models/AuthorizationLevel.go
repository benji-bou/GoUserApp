package models

type AuthorizationLevel byte

const (
	Public AuthorizationLevel = 1 << iota
	Basic
	Private
	Organisation
	Admin
	Root
)

type Authorizer interface {
	GetAuthorization() AuthorizationLevel
}

func (a AuthorizationLevel) MarshalJSON() ([]byte, error) {
	return []byte(a.Description()), nil
}

func (a AuthorizationLevel) Description() string {
	var result = ""
	switch {
	case a&Public != 0:
		result += "Public|"
	case a&Basic != 0:
		result += "Basic|"
	case a&Private != 0:
		result += "Private|"
	case a&Organisation != 0:
		result += "Organisation|"
	case a&Admin != 0:
		result += "Admin|"
	case a&Root != 0:
		result += "Root|"
	}
	if length := len(result); length > 0 {
		result = result[:length]
	}
	return result
}
