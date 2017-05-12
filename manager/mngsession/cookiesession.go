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
	manager := &CookieSessionManager{isSecure: isSecure, duration: duration, db: db, sessionType: reflect.TypeOf(session)}
	return manager
}

func (sm *CookieSessionManager) RemoveSession(user models.User) error {
	return sm.db.RemoveModel(dbm.M{"user._id": user.GetId()}, reflect.New(sm.sessionType))
}

func (sm *CookieSessionManager) ReadSession(c echo.Context, s models.Sessionizer) error {
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

func (sm *CookieSessionManager) WriteSession(c echo.Context, s models.Sessionizer) error {
	s.SetId(bson.NewObjectId())
	errs := sm.db.InsertModel(s)
	if errs != nil {
		return errs
	}
	cookie := sm.writeSessionCookie(s)
	c.SetCookie(cookie)
	return nil
}

func (sm *CookieSessionManager) writeSessionCookie(s models.Sessionizer) *http.Cookie {
	cookie := new(http.Cookie)
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
