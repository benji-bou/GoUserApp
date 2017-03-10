package mnguser

import (
	"errors"
	"goappuser/models"
)

var (
	ErrInvalidMail        = errors.New("invalid mail provided")
	ErrInvalidPassword    = errors.New("invalid password provided")
	ErrAlreadyRegister    = errors.New("mail already register")
	ErrAlreadyAuth        = errors.New("Already authenticate")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrNoSession          = errors.New("No session found")
)

//Manager interface to implements all feature to manage user
type UserManager interface {
	//Register register as a new user
	Register(models.User) error

	Update(models.User) error
	//IsExist check existence of the user
	IsExist(models.User) bool
	//ResetPassword user with specifics credentials
	ResetPassword(models.User, string) bool
	//GetByUniqueLogin retrieve a user using its UniqueLogin
	GetByEmail(Email string, user models.User) error
	//
	GetById(id string, user models.User) error
	//Authenticate
	Authenticate(username, password string, user models.User) (models.User, error)
	//Logout the current user
	Logout(user models.User) error
	//Add Friend
	AddFriend(user, friend models.User) error
	//User List
	UserList(login string, users interface{}) error
}
