package security

import (
	"github.com/labstack/echo"
	dbm "goappuser/database"
	"golang.org/x/crypto/bcrypt"
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

// func (ha *CustomHeaderAuth) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		s, _ := ha.ReadSessionUser(c)
// 		c.Set("Session", &s)
// 		next(c)
// 		return nil
// 	}
// }

// func (ha *CustomHeaderAuth) ReadSessionUser(c echo.Context) (*models.Session, error) {
// 	sessionId := c.Request().Header.Get(ha.tokenKey)
// 	s := &models.Session{}
// 	if errDB := ha.db.GetOneModel(dbm.M{"_id": sessionId}, s); errDB != nil {
// 		return nil, errDB
// 	}
// 	return s, nil
// }

// func (ha *CustomHeaderAuth) WriteSessionUser(c echo.Context, user models.User) error {
// 	if session, err := models.NewSession(); err != nil {
// 		log.Println("Session - CreateSession -", err)
// 		return err
// 	} else {
// 		session.User = user
// 		session.Id = bson.NewObjectId()
// 		errs := ha.db.InsertModel(session)
// 		c.Request().Header.Set(ha.tokenKey, session.Id.Hex())
// 		return errs
// 	}
// }
