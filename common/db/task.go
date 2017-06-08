import db

import "gopkg.in/mgo.v2/bson"

type Task struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Desc string
	Price uint32
	DeadLine int64
}
