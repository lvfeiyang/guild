package db

import "gopkg.in/mgo.v2/bson"

type Apply struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Desc string
	Price uint32
	ApplyTime int64
	DeadLine int64
}

const applyCName = "apply"
