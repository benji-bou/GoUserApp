package models

import (
	"fmt"
)

type AuthorizationLevel byte

const (
	Public AuthorizationLevel = 1 << iota
	Basic
	Private
	Organisation
	OrganisationAdmin
	Admin
	Root
)

type Authorizer interface {
	GetAuthorization() AuthorizationLevel
	AddAuthorization(newAuthlvl AuthorizationLevel)
}

func (a AuthorizationLevel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", a.Description(), "\"")), nil
}

func (a AuthorizationLevel) Description() string {
	var result = ""

	if a&Public != 0 {
		result += "Public|"
	}
	if a&Basic != 0 {
		result += "Basic|"
	}
	if a&Private != 0 {
		result += "Private|"
	}
	if a&Organisation != 0 {
		result += "Organisation|"
	}
	if a&OrganisationAdmin != 0 {
		result += "OrganisationAdmin|"
	}
	if a&Admin != 0 {
		result += "Admin|"
	}
	if a&Root != 0 {
		result += "Root|"
	}
	if length := len(result); length > 0 {
		result = result[:length-1]
	}
	return result
}
