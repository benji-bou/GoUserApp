package sessions

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/nu7hatch/gouuid"
	"goappuser/database"
	"log"
	"time"
)

var (
	manager *Manager = nil
)

var (
	ErrSessionManagerUnvailable = errors.New("The session manager hasn't beenn initialised with NewSessionManager")
)

type Manager struct {
	isSecure bool
	duration time.Time
	db       dbm.DatabaseQuerier
}

func NewManager(isSecure bool, duration time.Time, db dbm.DatabaseQuerier) *Manager {
	manager = &Manager{isSecure: isSecure, duration: duration, db: db}
	return manager
}

func (sm *Manager) CreateSession(c echo.Context, user interface{}) (Session, error) {
	if session, err := NewSession(c, user); err != nil {
		return session, err
	} else {
		errs := sm.db.InsertModel(session)
		return session, errs
	}
}

func (sm *Manager) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("MiddleWare Session")
		sessionId, err := readSessionCookie(c)
		if err != nil {
			log.Println("MiddleWare Session error readCookie : ", err)
			next(c)
			return err
		}
		s := &Session{}
		if errDB := sm.db.GetOneModel(dbm.M{"sessionId": sessionId.String()}, s); errDB != nil {
			log.Println("MiddleWare Session error GetOneModel : ", errDB)
			return errDB
		}
		c.Set("Session", *s)
		return nil
	}
}

type Session struct {
	id *uuid.UUID

	User   interface{}
	Values map[string]interface{}
}

func NewSession(c echo.Context, user interface{}) (Session, error) {
	log.Println("Create New Session", user)
	if manager == nil {
		log.Println("New Session erro", ErrSessionManagerUnvailable)
		return Session{}, ErrSessionManagerUnvailable
	}

	if uid, err := uuid.NewV4(); err != nil {
		log.Println("Create New Session uuid error", err)
		return Session{}, err
	} else {
		s := Session{id: uid,
			User:   user,
			Values: make(map[string]interface{})}
		writeCookie(c, s)
		return s, nil
	}
}

func writeCookie(c echo.Context, s Session) {
	cookie := new(echo.Cookie)
	cookie.SetName("sessionId")
	cookie.SetValue(s.id.String())
	cookie.SetExpires(manager.duration)
	c.SetCookie(cookie)
}

func readSessionCookie(c echo.Context) (*uuid.UUID, error) {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return nil, err
	}
	return uuid.ParseHex(cookie.Value())
}
