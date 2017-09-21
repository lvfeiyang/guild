package db

import "gopkg.in/mgo.v2/bson"
import "github.com/lvfeiyang/proxy/common/db"

type Guild struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Type      byte
	Introduce string
}

const guildCName = "guild"

func (g *Guild) GetById(id bson.ObjectId) error {
	return db.FindOneById(guildCName, id, g)
}
func (g *Guild) Save() error {
	g.Id = bson.NewObjectId()
	return db.Create(guildCName, g)
}
func (g *Guild) UpdateById() error {
	return db.UpdateOne(guildCName, g.Id, bson.M{"$set": bson.M{"introduce": g.Introduce}})
}
func FindAllGuilds() ([]Guild, error) {
	var gs []Guild
	// gs := make([]Guild, 20)
	err := db.FindMany(guildCName, bson.M{}, &gs, "")
	return gs, err
}
func DelGuildById(id bson.ObjectId) error {
	return db.DeleteOne(guildCName, id)
}
