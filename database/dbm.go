package dbm

import (
	"gotools"
	"log"
	"math/rand"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseQuerier interface {
	GetModels(mongoQuery M, resultInterface interface{}, limit int, skip int) (interface{}, error)
	GetOneModel(mongoQuery M, result interface{}) error
	InsertModel(model ...interface{}) []error
	IsExist(result interface{}) bool
}

type M map[string]interface{}

type MongoDatabaseSession struct {
	host     string
	port     string
	db_name  string
	username string
	password string

	Database *mgo.Database
}

type Modeler interface {
	SetName(name string)
	Name() string
}

func NewMongoDatabaseSession(host, port, db_name, username, password string) *MongoDatabaseSession {
	mongo := &MongoDatabaseSession{host, port, db_name, username, password, nil}
	log.Println("mongo loaded")
	return mongo
}

func configureMongoDatabaseSession(session *mgo.Session) {

}

func (db *MongoDatabaseSession) Connect() error {
	// log.Println("DB URL = " + db.host + ":" + db.port)
	//"mongodb://" + db.username + ":" + db.password +"@" +
	session, err := mgo.Dial(db.host + ":" + db.port)
	if err != nil {
		panic(err)
	}
	configureMongoDatabaseSession(session)
	db.Database = session.DB(db.db_name)
	return err
}

func (db *MongoDatabaseSession) Close() {
	db.Database.Session.Close()
}

//GetRandomOneModel get one model random in all documents
func (db *MongoDatabaseSession) GetRandomOneModel(model interface{}) error {
	collectionName := tools.GetInnerTypeName(model)
	collection := db.Database.C(collectionName)
	countElem, err := collection.Count()
	if err != nil {
		return err
	}
	elemNum := rand.Intn(countElem)
	collection.Find(bson.M{}).Skip(elemNum).One(model)
	return nil
}

//GetModels retrieves all the data from mongoDB
//mongoQuery is the query from MongoDB query
//resultInterface is a slice representing the model required, it will be fill with the result of the query
//limit of result if limit < 0 no limit used
//skip corresponding the number elements to skip
//return an err if soimething bad appened
func (db *MongoDatabaseSession) GetModels(mongoQuery M, resultInterface interface{}, limit int, skip int) (interface{}, error) {

	collectionName := tools.GetInnerTypeName(resultInterface)
	collection := db.Database.C(collectionName)
	result := tools.CreatePtrToSliceFromInterface(resultInterface)
	var err error = nil
	switch {
	case limit <= 0 && skip <= 0:
		err = collection.Find(bson.M(mongoQuery)).All(result)
	case limit > 0 && skip <= 0:
		err = collection.Find(bson.M(mongoQuery)).Limit(limit).All(result)
	case limit <= 0 && skip > 0:
		err = collection.Find(bson.M(mongoQuery)).Skip(skip).All(result)
	case limit > 0 && skip > 0:
		err = collection.Find(bson.M(mongoQuery)).Skip(skip).Limit(limit).All(result)
	}
	resultInterface = tools.Dereference(result)
	return resultInterface, err
}

func (db *MongoDatabaseSession) GetOneModel(mongoQuery M, result interface{}) error {
	collectionName := tools.GetInnerTypeName(result)
	log.Println(collectionName)
	collection := db.Database.C(collectionName)
	err := collection.Find(bson.M(mongoQuery)).One(result)
	return err
}

// func (db *MongoDatabaseSession) GetModel(collectionName string, mongoQuery M, result interface{}) error {
// 	collection := db.Database.C(collectionName)
// 	err := collection.Find(bson.M(mongoQuery)).All(result)
// 	return err
// }

func (db *MongoDatabaseSession) InsertModel(model ...interface{}) []error {
	sortedModel := tools.SortArrayByType(model)
	err := make([]error, 0)
	for collectionName, models := range sortedModel {
		collection := db.Database.C(collectionName)
		errTmp := collection.Insert(models...)
		if errTmp != nil {
			err = append(err, errTmp)
		}
	}
	if len(err) == 0 {
		err = nil
	}
	return err
}

func (db *MongoDatabaseSession) IsExist(result interface{}) bool {
	newResult := tools.Zero(result)
	queryMap, _ := tools.Map(result)
	if len(queryMap) == 0 {
		return false
	}
	db.GetOneModel(queryMap, &newResult)
	return tools.NotEmpty(newResult)
}
