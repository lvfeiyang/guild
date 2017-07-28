package db

import "gopkg.in/mgo.v2/bson"

type Member struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Mobile  string
	Ability string
	Name string
	GuildId string
	Role byte
	Accounts []string
}

const memberCName = "member"

const (
	RoleMaster = 1<<iota
	RoleAdmin
)

func (m *Member) Save() error {
	m.Id = bson.NewObjectId()
	return Create(memberCName, m)
}
func (m *Member) UpdateById() error {
	return UpdateOne(memberCName, m.Id, bson.M{"$set": bson.M{"mobile": m.Mobile}})
}
func FindAllMembers(gId string) ([]Member, error) {
	var ms []Member
	err := FindMany(memberCName, bson.M{"guildid": gId}, &ms)
	return ms, err
}
func (m *Member) GetById(id bson.ObjectId) error {
	return FindOneById(memberCName, id, m)
}
func DelMemberById(id bson.ObjectId) error {
	return DeleteOne(memberCName, id)
}
func DelMembersByGId(gId string) error {
	return DeleteMany(memberCName, bson.M{"guildid": gId})
}
