package mngsession

import (
	"github.com/labstack/echo"
	dbm "goappuser/database"
	"goappuser/models"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"reflect"
	"time"
)

type CookieSessionManager struct {
	isSecure    bool
	duration    time.Time
	db          dbm.DatabaseQuerier
	sessionType reflect.Type
}

func NewCookieSessionManager(isSecure bool, duration time.Time, db dbm.DatabaseQuerier, session models.Sessionizer) *CookieSessionManager {
	log.Println("type of sessionizer = ", reflect.TypeOf(session))
	manager := &CookieSessionManager{isSecure: isSecure, duration: duration, db: db, sessionType: reflect.TypeOf(session)}
	return manager
}

func (sm *CookieSessionManager) RemoveSession(user models.User) error {
	return sm.db.RemoveModel(dbm.M{"user._id": user.GetId()}, reflect.New(sm.sessionType))
}

func (sm *CookieSessionManager) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		inter := reflect.New(sm.sessionType.Elem()).Interface()

		log.Println("Middleware - CookieSession", reflect.TypeOf(inter))

		if session, isOK := inter.(models.Sessionizer); isOK == true {
			log.Println("getSession Ok ", reflect.TypeOf(session))

			err := sm.ReadSessionUser(c, session)
			if err != nil && err != ErrNoSessionFound {
				log.Println("session error Middleware", err)
			}
		} else {
			log.Println("session error Middleware --> Problem getting a sessionizer")
		}
		next(c)
		return nil
	}
}

func (sm *CookieSessionManager) ReadSessionUser(c echo.Context, s models.Sessionizer) error {
	bsonId, err := sm.readSessionCookie(c)
	if err != nil {
		return ErrNoSessionFound
	}
	log.Println("id found", bsonId)
	if errDB := sm.db.GetOneModel(dbm.M{"_id": bsonId}, s); errDB != nil {
		return errDB
	}
	c.Set("Session", s)
	return nil
}

func (sm *CookieSessionManager) WriteSessionUser(c echo.Context, user models.User) error {
	inter := reflect.New(sm.sessionType.Elem()).Interface()

	// log.Println("WriteSessionUser", inter.Name())

	if session, isOK := inter.(models.Sessionizer); isOK == true {
		session.SetUser(user)
		session.SetId(bson.NewObjectId())
		errs := sm.db.InsertModel(session)
		if errs != nil {
			return errs
		}
		cookie := sm.writeSessionCookie(session)
		c.SetCookie(cookie)
		return nil

	} else {
		log.Println("session error Middleware", ErrSessionNotSessionizer)
		return ErrSessionNotSessionizer
	}
}

func (sm *CookieSessionManager) writeSessionCookie(s models.Sessionizer) *http.Cookie {
	cookie := new(http.Cookie)
	//	cookie.SetSecure(true)
	// cookie.SetHTTPOnly(true)
	cookie.Name = "sessionId"
	cookie.Value = s.GetId().Hex()
	cookie.Path = "/"
	cookie.Expires = sm.duration
	return cookie
}

func (sm *CookieSessionManager) readSessionCookie(c echo.Context) (bson.ObjectId, error) {
	cookie, err := c.Cookie("sessionId")
	if err != nil {
		log.Println("Read Cookie err", err)
		return bson.ObjectId(""), err
	}
	strId := cookie.Value
	log.Println("id of cookie --> ", strId)
	return bson.ObjectIdHex(strId), nil
}
