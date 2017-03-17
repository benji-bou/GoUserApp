package security

import (
	"github.com/labstack/echo"
	dbm "goappuser/database"
	"goappuser/manager/mngsession"
	"goappuser/models"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func NewCustomHeaderAuth(db dbm.DatabaseQuerier, usernameKey, passwordKey, tokenKey string) *CustomHeaderAuth {
	return &CustomHeaderAuth{db: db, usernameKey: usernameKey, passwordKey: passwordKey, tokenKey: tokenKey}
}

type CustomHeaderAuth struct {
	db          dbm.DatabaseQuerier
	usernameKey string
	passwordKey string
	tokenKey    string
}

func (ha *CustomHeaderAuth) GetCredentials(c echo.Context) (string, string, error) {
	return c.Request().Header.Get(ha.usernameKey), c.Request().Header.Get(ha.passwordKey), nil
}
func (ha *CustomHeaderAuth) Hash(clearpassword []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(clearpassword, 0)
}
func (ha *CustomHeaderAuth) Compare(clearPassword, realmPassword []byte) bool {
	errcmp := bcrypt.CompareHashAndPassword(realmPassword, clearPassword)
	return errcmp == nil
}

func (ha *CustomHeaderAuth) ReadSession(c echo.Context, s models.Sessionizer) error {
	token := c.Request().Header.Get(ha.tokenKey)
	if token == "" {
		log.Println("In header request didn't found", ha.tokenKey, "key", c.Request().Header)
		return mngsession.ErrNoSessionFound
	}
	return ha.db.GetOneModel(dbm.M{"_id": bson.ObjectIdHex(token)}, s)
}

func (ha *CustomHeaderAuth) WriteSession(c echo.Context, s models.Sessionizer) error {
	sessionId := bson.NewObjectId()
	c.Response().Header().Set(ha.tokenKey, sessionId.Hex())
	s.SetId(sessionId)
	return ha.db.InsertModel(s)
}
