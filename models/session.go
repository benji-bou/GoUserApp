package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Sessionizer interface {
	GetId() bson.ObjectId
	SetId(id bson.ObjectId)
	GetUser() User
	SetUser(user User)
}

type Session struct {
	Id bson.ObjectId `bson:"_id" json:"id"`
}

func NewSession() (Session, error) {
	s := Session{
		Id: bson.NewObjectId(),
	}
	return s, nil
}

func (s *Session) SetId(id bson.ObjectId) {
	s.Id = id
}

func (s *Session) GetId() bson.ObjectId {
	return s.Id
}
