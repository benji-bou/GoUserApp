package models

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	ErrUserFriendInvalid  = errors.New("Users is not a valid friend")
	ErrUserAlreadyFriend  = errors.New("Users is already in the friend list")
	ErrUserFriendNotFound = errors.New("Users not found in the friend list")
)

//NewUser create a basic user with the mandatory parameters for each users
func NewUserDefaultExtended(Email, password string, role AuthorizationLevel) *UserDefaultExtended {

	return &UserDefaultExtended{DateCreate: time.Now(), UserDefault: UserDefault{Email: Email, Password: []byte(password), Role: role, Friends: make([]UserDefault, 0, 0)}}
}

func NewUserDefault(user User) *UserDefault {

	return &UserDefault{Id: user.GetId(), Email: user.GetEmail(), Password: user.GetPassword(), Role: user.GetAuthorization()}
}

//User Represent a basic user

type User interface {
	Authorizer
	SetId(id bson.ObjectId)
	GetId() bson.ObjectId
	GetEmail() string
	SetEmail(Email string)
	GetPassword() []byte
	SetPassword(pass []byte)
	GetFriends() []User
	AddFriend(user User) error
	RemoveFriend(user User) error
}

//User Represent a basic user

type UserDefault struct {
	Id       bson.ObjectId      `bson:"_id" json:"id"`
	Password []byte             `bson:"password" json:"-"`
	Email    string             `bson:"email" json:"email"`
	Role     AuthorizationLevel `bson:"roles" json:"roles"`
	Friends  []UserDefault      `bson:"friends" json:"friends"`
}

type UserDefaultExtended struct {
	UserDefault        `bson:"credentials,inline"`
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

func (u *UserDefault) GetEmail() string {
	return u.Email
}

func (u *UserDefault) SetEmail(Email string) {
	u.Email = Email

}
func (u *UserDefault) GetPassword() []byte {
	return u.Password
}

func (u *UserDefault) SetPassword(pass []byte) {
	u.Password = pass
}

func (u *UserDefault) GetAuthorization() AuthorizationLevel {
	return u.Role
}

func (u *UserDefault) AddAuthorization(newAuthlvl AuthorizationLevel) {
	u.Role = u.Role | newAuthlvl
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

func (u *UserDefault) UnmarshalJSON(data []byte) error {
	type aliasUserDefault UserDefault
	usr := &struct {
		aliasUserDefault
		Pass string `json:"password"`
	}{}

	if err := json.Unmarshal(data, usr); err != nil {
		return err
	}
	*u = UserDefault{Email: usr.Email, Password: []byte(usr.Pass), Role: usr.Role, Id: usr.Id, Friends: usr.Friends}

	return nil
}
