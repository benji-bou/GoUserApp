package models

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"goappuser/database"
	"goappuser/security"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
	"gotools"
	"log"
	"time"
)

var (
	ErrInvalidMail        = errors.New("invalid mail provided")
	ErrInvalidPassword    = errors.New("invalid password provided")
	ErrAlreadyRegister    = errors.New("mail already register")
	ErrAlreadyAuth        = errors.New("Already authenticate")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrNoSession          = errors.New("No session found")
	ErrUserFriendInvalid  = errors.New("Users is not a valid friend")
	ErrUserAlreadyFriend  = errors.New("Users is already in the friend list")
	ErrUserFriendNotFound = errors.New("Users not found in the friend list")
)

//Manager interface to implements all feature to manage user
type UserManager interface {
	//Register register as a new user
	Register(User) error

	Update(User) error
	//IsExist check existence of the user
	IsExist(User) bool
	//ResetPassword user with specifics credentials
	ResetPassword(User, string) bool
	//GetByUniqueLogin retrieve a user using its UniqueLogin
	GetByUniqueLogin(UniqueLogin string, user User) error
	//
	GetById(id string, user User) error
	//Authenticate
	Authenticate(c echo.Context, user User) (User, error)
	//Logout the current user
	Logout(user User) error
	//Add Friend
	AddFriend(user, friend User) error
	//User List
	UserList(login string, user User) (interface{}, error)
}

//NewUser create a basic user with the mandatory parameters for each users
func NewUserDefaultExtended(UniqueLogin, password string) *UserDefaultExtended {
	log.Println("New Password", password)
	return &UserDefaultExtended{UserDefault: UserDefault{UniqueLogin: UniqueLogin, Password: []byte(password), Role: "user", Friends: make([]UserDefault, 0, 0)}}
}

func NewUserDefault(user User) *UserDefault {

	return &UserDefault{Id: user.GetId(), UniqueLogin: user.GetUniqueLogin(), Password: user.GetPassword(), Role: user.GetRole()}
}

//User Represent a basic user

//TODO: Change User to an interface
type User interface {
	SetId(id bson.ObjectId)
	GetId() bson.ObjectId
	GetUniqueLogin() string
	SetUniqueLogin(UniqueLogin string)
	GetPassword() []byte
	SetPassword(pass []byte)
	GetRole() string
	SetRole(role string)
	GetFriends() []User
	AddFriend(user User) error
	RemoveFriend(user User) error
}

//User Represent a basic user

type UserDefault struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Password    []byte        `bson:"password" json:"-"`
	UniqueLogin string        `bson:"uniqueLogin" json:"uniqueLogin"`
	Role        string        `bson:"role" json:"-"`
	Friends     []UserDefault `bson:"friends" json:"friends"`
}

type UserDefaultExtended struct {
	UserDefault        `bson:"credentials,inline" json:"credentials,inline"`
	Email              string    `bson:"email" json:"email"`
	Name               string    `bson:"name" json:"name"`
	Surname            string    `bson:"surname" json:"surname"`
	Pseudo             string    `bson:"pseudo" json:"pseudo"`
	DateCreate         time.Time `bson:"created" json:"created"`
	DateLastConnection time.Time `bson:"lastconnection" json:"lastconnection,omitempty"`
	BirthDate          time.Time `bson:"birthdate" json:"birthdate,omitempty"`
}

func (u *UserDefault) SetId(id bson.ObjectId) {
	u.Id = id
}

func (u *UserDefault) GetId() bson.ObjectId {
	return u.Id
}

func (u *UserDefault) GetUniqueLogin() string {
	return u.UniqueLogin
}

func (u *UserDefault) SetUniqueLogin(UniqueLogin string) {
	u.UniqueLogin = UniqueLogin

}
func (u *UserDefault) GetPassword() []byte {
	return u.Password
}

func (u *UserDefault) SetPassword(pass []byte) {
	u.Password = pass
}

func (u *UserDefault) GetRole() string {
	return u.Role
}

func (u *UserDefault) SetRole(role string) {
	u.Role = role
}

func (u *UserDefault) GetFriends() []User {
	count := len(u.Friends)
	res := make([]User, count, count)
	for i, s := range u.Friends {
		res[i] = &s
	}
	return res
}

func (u *UserDefault) AddFriend(user User) error {
	if u.GetId() == user.GetId() {
		return ErrUserFriendInvalid
	}
	for _, fr := range u.Friends {
		if fr.GetId() == user.GetId() {
			return ErrUserAlreadyFriend
		}
	}
	u.Friends = append(u.Friends, *NewUserDefault(user))
	return nil
}

func (u *UserDefault) RemoveFriend(user User) error {
	for index, fr := range u.Friends {
		if fr.GetId() == user.GetId() {
			u.Friends = append(u.Friends[:index], u.Friends[index+1:]...)
			return nil
		}
	}
	return ErrUserFriendNotFound
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
func (m *DBUserManage) Register(user User) error {
	if user.GetUniqueLogin() == "" {
		return ErrInvalidMail
	}

	if user.GetPassword() == nil {
		log.Println(string(user.GetPassword()))
		return ErrInvalidPassword
	}

	if m.IsExist(user) {
		return ErrAlreadyRegister
	}

	pass, errHash := m.auth.Hash(user.GetPassword())
	if errHash != nil {
		return errHash
	}
	user.SetId(bson.NewObjectId())
	user.SetPassword(pass)
	log.Println("insert user", user)
	if errInsert := m.db.InsertModel(user); errInsert != nil {
		log.Println("error insert", errInsert, " user: ", user.GetUniqueLogin())
		return errInsert
	}
	log.Println("insert OK")
	return nil
}

func (m *DBUserManage) Update(user User) error {
	log.Println(user.GetFriends()[0].GetId())
	return m.db.UpdateModelId(user.GetId(), user)
}

//IsExist check existence of the user
func (m *DBUserManage) IsExist(user User) bool {
	u := &UserDefault{}
	if err := m.GetByUniqueLogin(user.GetUniqueLogin(), u); err == nil {
		log.Println("IsExist user ", u)
		return tools.NotEmpty(u)
	} else if err == mgo.ErrNotFound {
		return false
	}
	return false
}

//ResetPassword user with specifics credentials
func (m *DBUserManage) ResetPassword(user User, newPassword string) bool {
	return false
}

//GetByUniqueLogin retrieve a user using its UniqueLogin
func (m *DBUserManage) GetByUniqueLogin(UniqueLogin string, user User) error {
	if err := m.db.GetOneModel(dbm.M{"UniqueLogin": UniqueLogin}, user); err != nil {
		return err
	}
	return nil
}

//GetById retrieve a user using its id
func (m *DBUserManage) GetById(id string, user User) error {
	if bson.IsObjectIdHex(id) == false {
		return ErrUserNotFound
	}
	if err := m.db.GetOneModel(dbm.M{"_id": bson.ObjectIdHex(id)}, user); err != nil {
		return err
	}
	return nil
}

//Authenticate log the user
func (m *DBUserManage) Authenticate(c echo.Context, user User) (User, error) {
	if session, isOk := (c).Get("Session").(Session); isOk {
		if err := m.GetByUniqueLogin(session.User.GetUniqueLogin(), user); err != nil {
			return nil, ErrUserNotFound
		}
		return user, ErrAlreadyAuth
	}
	username, password, err := m.auth.GetCredentials(c)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Failed to retrieve credentials from request: ", err))
	}
	err = m.GetByUniqueLogin(username, user)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if ok := m.auth.Compare([]byte(password), user.GetPassword()); ok == true {
		if _, cookie, err := m.sessionManager.CreateSession(user); err == nil {
			(c).SetCookie(cookie)
		} else {
			log.Println(err)
			return user, err
		}
		return user, nil
	}
	log.Println(ErrInvalidCredentials.Error(), username, password)
	return nil, ErrInvalidCredentials
}

func (m *DBUserManage) Logout(user User) error {
	return m.sessionManager.RemoveSession(user)
}

func (m *DBUserManage) AddFriend(user, friend User) error {
	if err := user.AddFriend(friend); err != nil {
		log.Println("arrFriend disk error")
		return err
	} else {
		return m.Update(user)
	}
}

func (m *DBUserManage) UserList(login string, user User) (interface{}, error) {
	UniqueLogin := fmt.Sprintf(".*%s.*", login)

	//dbm.M{"$regex":
	return m.db.GetModels(dbm.M{"UniqueLogin": bson.RegEx{UniqueLogin, ""}}, &user, 20, 0)
}

func (m *DBUserManage) cleanSession(c echo.Context) error {
	if _, isOk := c.Get("Session").(Session); isOk {
		return ErrNoSession
	}
	return nil
	//TODO: use m.db.Remove Model to remove the session
	//	m.db.
}
