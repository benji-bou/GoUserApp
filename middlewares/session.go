package middlewares

import (
	"github.com/labstack/echo"
	"github.com/nu7hatch/gouuid"
	"goappuser/database"
	"goappuser/user"
)

type SessionManagerUnvailable error

var (
	manager SessionManager = nil
)

type SessionManager struct {
	isSecure bool
	duration int
	db       dbm.DatabaseQuerier
}

func NewSessionManager(isSecure bool, duration int, db dbm.DatabaseQuerier) SessionManager {
	manager = SessionManager{IsSecure: isSecure, duration: duration, db: db}
	return manager
}

type Session struct {
	id uuid.UUID

	User   *user.User
	Values map[string]interface{}
}

func NewSession(c echo.Context, user *user.User, db dbm.DatabaseQuerier) (error, *Session) {
	if manager == nil {
		return SessionManagerUnvailable, nil
	}

	s := &Session{id: uuid.NewV5(uuid.NamespaceURL, []byte(user.Email)),
		config: SessionConfig{IsSecure: true, duration: time.Now().Add(24 * time.Hour), db: db},
		User:   user,
		Values: make(map[string]interface{})}
	s.writeCookie(c)
	return nil, s
}

func (s *Session) writeCookie(c echo.Context) {
	cookie := new(echo.Cookie)
	cookie.SetName("sessionId")
	cookie.SetValue(s.id.String())
	cookie.SetExpires(manager.duration)
	c.SetCookie(cookie)
}

func readSessionCookie(c echo.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		return nil, err
	}
	return cookie.Value(), nil
}

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if manager == nil {
			return SessionManagerUnvailable
		}
		sessionId, err := readSessionCookie(c)
		if err != nil {
			return err
		}
		s := &Session{}
		if errDB := manager.db.GetOneModel(dbm.M{"sessionId": sessionId.String()}, s); errDB != nil {
			return errDB
		}
		c.Set("Session", s)
		next(c)
	}
}
