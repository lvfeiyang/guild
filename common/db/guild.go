import db

import "gopkg.in/mgo.v2/bson"

type Guild struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string
	Type byte
	Introduce string
}
