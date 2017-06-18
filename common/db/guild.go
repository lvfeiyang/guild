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
func (g *Guild) Save() error {
	g.Id = bson.NewObjectId()
	return Create(guildCName, g)
}
func (g *Guild) UpdateById() error {
	return UpdateOne(guildCName, g.Id, bson.M{"$set": bson.M{"introduce": g.Introduce}})
}
func FindAllGuilds() ([]Guild, error) {
	var gs []Guild
	// gs := make([]Guild, 20)
	err := FindMany(guildCName, bson.M{}, &gs)
	return gs, err
}
