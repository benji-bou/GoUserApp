package models

import (
	"errors"
	"github.com/labstack/echo"
	"goappuser/database"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var (
	manager *SessionManager = nil
)

var (
	ErrSessionManagerUnvailable = errors.New("The session manager hasn't beenn initialised with NewSessionManager")
)

type SessionManager struct {
	isSecure bool
	duration time.Time
	db       dbm.DatabaseQuerier
}

func NewSessionManager(isSecure bool, duration time.Time, db dbm.DatabaseQuerier) *SessionManager {
	manager = &SessionManager{isSecure: isSecure, duration: duration, db: db}
	return manager
}

func (sm *SessionManager) CreateSession(user *User) (Session, *echo.Cookie, error) {
	if session, err := NewSession(user); err != nil {
		log.Println("error: ", err)
		return session, nil, err
	} else {
		errs := sm.db.InsertModel(session)
		return session, writeSessionCookie(session), errs
	}
}

func (sm *SessionManager) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionId, err := readSessionCookie(c)
		if err != nil {
			next(c)
			return err
		}
		s := &Session{}
		if errDB := sm.db.GetOneModel(dbm.M{"_id": sessionId}, s); errDB != nil {
			next(c)
			return errDB
		}
		c.Set("Session", *s)
		next(c)
		return nil
	}
}

type Session struct {
	Id     bson.ObjectId          `bson:"_id" json:"id"`
	User   *User                  `bson:"user" json:"user"`
	Values map[string]interface{} `bson:"values" json:"values"`
}

func NewSession(user *User) (Session, error) {
	if manager == nil {
		return Session{}, ErrSessionManagerUnvailable
	}
	s := Session{
		Id:     bson.NewObjectId(),
		User:   user,
		Values: make(map[string]interface{})}
	return s, nil
}

func writeSessionCookie(s Session) *echo.Cookie {
	cookie := new(echo.Cookie)
	cookie.SetName("sessionId")
	cookie.SetValue(s.Id.Hex())
	cookie.SetExpires(manager.duration)
	return cookie
}

func readSessionCookie(c echo.Context) (bson.ObjectId, error) {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return bson.ObjectId(""), err
	}
	strId := cookie.Value()
	return bson.ObjectIdHex(strId), nil
}
