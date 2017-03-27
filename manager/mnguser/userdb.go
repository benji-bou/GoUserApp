package mnguser

import (
	"fmt"
	dbm "goappuser/database"
	"goappuser/models"
	"goappuser/security"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gotools/reflectutil"
	"log"
	"reflect"
	"time"
)

//NewDBUserManage create a db manager user
func NewDBUserManage(db dbm.DatabaseQuerier, auth security.AuthenticationProcesser) *DBUserManage {
	return &DBUserManage{db: db, auth: auth}
}

//DBUserManage respect Manager Interface using MGO (MongoDB driver)
type DBUserManage struct {
	db   dbm.DatabaseQuerier
	auth security.AuthenticationProcesser
}

//Register register as a new user
func (m *DBUserManage) Register(user models.User) error {
	log.Println("starting register ", user.GetEmail())
	if user.GetEmail() == "" {
		log.Println("GetEmail is empty return ", ErrInvalidMail.Error())
		return ErrInvalidMail
	}

	if user.GetPassword() == nil {
		log.Println("GetPassword is nil return ", ErrInvalidPassword.Error())
		log.Println(string(user.GetPassword()))
		return ErrInvalidPassword
	}

	if m.IsExist(user) {
		log.Println("user exist return ", ErrAlreadyRegister.Error())
		return ErrAlreadyRegister
	}
	pass, errHash := m.auth.Hash(user.GetPassword())
	if errHash != nil {
		log.Println("failed to hash with error", errHash.Error())
		return errHash
	}
	user.SetId(bson.NewObjectId())
	user.SetPassword(pass)
	if userExtended, isOk := user.(models.DateConnectionTracker); isOk {
		userExtended.SetCreationDate(time.Now())
		userExtended.SetNewConnectionDate(time.Now())
	}
	log.Println("insert user", user)
	if errInsert := m.db.InsertModel(user); errInsert != nil {
		log.Println("error insert", errInsert, " user: ", user.GetEmail())
		return errInsert
	}
	log.Println("insert OK")
	return nil
}

func (m *DBUserManage) Update(user models.User) error {
	return m.db.UpdateModelId(user.GetId(), user)
}

//IsExist check existence of the user
func (m *DBUserManage) IsExist(user models.User) bool {
	u := &models.UserDefault{}
	if err := m.GetByEmail(user.GetEmail(), u); err == nil {
		log.Println("IsExist user ", u)
		return reflectutil.NotEmpty(u)
	} else if err == mgo.ErrNotFound {
		return false
	}
	return false
}

//ResetPassword user with specifics credentials
func (m *DBUserManage) ResetPassword(user models.User, newPassword string) bool {
	return false
}

//GetByUniqueLogin retrieve a user using its UniqueLogin
func (m *DBUserManage) GetByEmail(email string, user models.User) error {
	log.Println("retrieve from email", email)
	if err := m.db.GetOneModel(dbm.M{"email": email}, user); err != nil {
		log.Println("failed retrieve from email", err, reflect.TypeOf(user))
		return err
	}
	return nil
}

//GetById retrieve a user using its id
func (m *DBUserManage) GetById(id string, user models.User) error {
	if bson.IsObjectIdHex(id) == false {
		return ErrUserNotFound
	}
	if err := m.db.GetOneModel(dbm.M{"_id": bson.ObjectIdHex(id)}, user); err != nil {
		return err
	}
	return nil
}

//Authenticate log the user with username and password. Try to retrieve models.User type passed in param
func (m *DBUserManage) Authenticate(username, password string, user models.User) error {
	err := m.GetByEmail(username, user)
	if err != nil {
		return ErrUserNotFound
	}
	if userExtended, isOk := user.(models.DateConnectionTracker); isOk {
		userExtended.SetNewConnectionDate(time.Now())
		m.db.UpdateModelId(user.GetId(), userExtended)
	}
	if ok := m.auth.Compare([]byte(password), user.GetPassword()); ok == true {
		return nil
	}

	log.Println(ErrInvalidCredentials.Error(), username, password)
	return ErrInvalidCredentials
}

func (m *DBUserManage) AddFriend(user, friend models.User) error {
	if err := user.AddFriend(friend); err != nil {
		log.Println("arrFriend disk error")
		return err
	} else {
		return m.Update(user)
	}
}

func (m *DBUserManage) Logout(user models.User) error {
	return nil
}

///UserList make the request to the DB and fill users
///users paramater must be a pointer to slice of user where user type used as CollectionName
func (m *DBUserManage) UserList(login string, users interface{}) error {
	UniqueLogin := fmt.Sprintf(".*%s.*", login)
	err := m.db.GetModels(dbm.M{"email": bson.RegEx{UniqueLogin, ""}}, users, 20, 0)
	return err
}
