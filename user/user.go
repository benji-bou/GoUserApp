package user

import (
	"errors"
	"fmt"
	"goappuser/database"
	"goappuser/security"
	"log"
	"net/http"
	"time"
)

//Manager interface to implements all feature to manage user
type Manager interface {
	//Register register as a new user
	Register(*User) error
	//IsExist check existence of the user
	IsExist(*User) bool
	//ResetPassword user with specifics credentials
	ResetPassword(*User, string) bool
	//GetByEmail retrieve a user using its email
	GetByEmail(string) (*User, error)
	//Authenticate
	Authenticate(r *http.Request) (*User, error)
}

//NewUser create a basic user with the mandatory parameters for each users
func NewUser(email, password string) *User {
	return &User{Email: email, Password: []byte(password), Role: "user"}
}

//User Represent a basic user
type User struct {
	Name               string      `bson:"name" json:"name"`
	Surname            string      `bson:"surname" json:"surname"`
	Pseudo             string      `bson:"pseudo" json:"pseudo"`
	Password           []byte      `bson:"password" json:"-"`
	Email              string      `bson:"email" json:"email"`
	DateCreate         time.Time   `bson:"created" json:"created"`
	DateLastConnection time.Time   `bson:"lastconnection" json:"lastconnection,omitempty"`
	BirthDate          time.Time   `bson:"birthdate" json:"birthdate,omitempty"`
	AdditionalInfos    interface{} `bson:"infos" json:"infos,omitempty"`
	Role               string      `bson:"role" json:"-"`
}

//NewDBUserManager create a db manager user
func NewDBUserManager(db dbm.DatabaseQuerier, auth security.AuthenticationProcesser) *DBUserManager {
	return &DBUserManager{db: db, auth: auth}
}

//DBUserManager respect Manager Interface using MGO (MongoDB driver)
type DBUserManager struct {
	db   dbm.DatabaseQuerier
	auth security.AuthenticationProcesser
}

//Register register as a new user
func (m *DBUserManager) Register(user *User) error {
	if !m.IsExist(user) {
		return errors.New("mail already register")
	}
	pass, errHash := m.auth.Hash(user.Password)
	if errHash != nil {
		return errHash
	}
	user.Password = pass
	if errInsert := m.db.InsertModel(user); len(errInsert) > 0 {
		return errInsert[0]
	}
	return nil
}

//IsExist check existence of the user
func (m *DBUserManager) IsExist(user *User) bool {
	return m.db.IsExist(user)
}

//ResetPassword user with specifics credentials
func (m *DBUserManager) ResetPassword(user *User, newPassword string) bool {
	return false
}

//GetByEmail retrieve a user using its email
func (m *DBUserManager) GetByEmail(email string) (*User, error) {
	user := &User{}
	if err := m.db.GetOneModel(dbm.M{"email": email}, user); err != nil {
		return nil, err
	}
	return user, nil
}

//Authenticate log the user
func (m *DBUserManager) Authenticate(c echo.Context) (*User, error) {
	username, password, err := m.auth.GetCredentials(c)
	if err != nil {
		log.Println("Login Error :", err)
		return nil, errors.New(fmt.Sprint("Failed to retrieve credentials from request: ", err))
	}

	user, err := m.GetByEmail(username)
	if err != nil {
		log.Println("Error logged in:", err)
		return nil, errors.New("User not found")
	}
	if ok := m.auth.Compare([]byte(password), user.Password); ok == true {

		return user, nil
	}
	return nil, errors.New("Invalid credentials")
}
