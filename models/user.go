package models

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"goappuser/database"
	"goappuser/security"

	"gopkg.in/mgo.v2"
	"gotools"
	"log"
	"time"
)

//Manager interface to implements all feature to manage user
type UserManager interface {
	//Register register as a new user
	Register(*User) error
	//IsExist check existence of the user
	IsExist(*User) bool
	//ResetPassword user with specifics credentials
	ResetPassword(*User, string) bool
	//GetByEmail retrieve a user using its email
	GetByEmail(string) (*User, error)
	//Authenticate
	Authenticate(c *echo.Context) (*User, error)
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

//NewDBUserManage create a db manager user
func NewDBUserManage(db dbm.DatabaseQuerier, auth security.AuthenticationProcesser, sm *SessionManager) *DBUserManage {
	return &DBUserManage{db: db, auth: auth, sessionManager: sm}
}

//DBUserManage respect Manager Interface using MGO (MongoDB driver)
type DBUserManage struct {
	db             dbm.DatabaseQuerier
	auth           security.AuthenticationProcesser
	sessionManager *SessionManager
}

//Register register as a new user
func (m *DBUserManage) Register(user *User) error {
	if m.IsExist(user) {
		return errors.New("mail already register")
	}

	pass, errHash := m.auth.Hash(user.Password)
	if errHash != nil {
		return errHash
	}
	user.Password = pass
	if errInsert := m.db.InsertModel(user); errInsert != nil {
		log.Println("error insert", errInsert, " user: ", user.Email)
		return errInsert
	}
	log.Println("insert OK")
	return nil
}

//IsExist check existence of the user
func (m *DBUserManage) IsExist(user *User) bool {
	if u, err := m.GetByEmail(user.Email); err == nil {
		log.Println("IsExist user ", u)
		return tools.NotEmpty(u)
	} else if err == mgo.ErrNotFound {
		return false
	}
	return false
}

//ResetPassword user with specifics credentials
func (m *DBUserManage) ResetPassword(user *User, newPassword string) bool {
	return false
}

//GetByEmail retrieve a user using its email
func (m *DBUserManage) GetByEmail(email string) (*User, error) {
	user := &User{}
	if err := m.db.GetOneModel(dbm.M{"email": email}, user); err != nil {
		return nil, err
	}
	return user, nil
}

//Authenticate log the user
func (m *DBUserManage) Authenticate(c *echo.Context) (*User, error) {
	if session, isOk := (*c).Get("Session").(Session); isOk {
		return session.User, errors.New("Already authenticate")
	}
	username, password, err := m.auth.GetCredentials(*c)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to retrieve credentials from request: ", err))
	}

	user, err := m.GetByEmail(username)
	if err != nil {
		return nil, errors.New("User not found")
	}
	if ok := m.auth.Compare([]byte(password), user.Password); ok == true {

		if _, cookie, err := m.sessionManager.CreateSession(user); err == nil {
			(*c).SetCookie(cookie)
		}
		return user, nil
	}
	return nil, errors.New("Invalid credentials")
}
