package db

import "gopkg.in/mgo.v2/bson"

type Member struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Mobile string
	Guild bson.ObjectId
}

const memberCName = "member"
