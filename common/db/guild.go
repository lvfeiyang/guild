package db

import "gopkg.in/mgo.v2/bson"

type Guild struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string
	Type byte
	Introduce string
}

const guildCName = "guild"

func (g *Guild) GetById(id bson.ObjectId) error {
	return FindOneById(guildCName, id, g)
}
func (ac *Guild) Save() error {
	ac.Id = bson.NewObjectId()
	return Create(guildCName, ac)
}
