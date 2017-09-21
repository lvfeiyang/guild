package db

import (
	"errors"
	"github.com/lvfeiyang/proxy/common/db"
	"github.com/lvfeiyang/proxy/common/session"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Member struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Mobile   string        //TODO: same account
	Ability  string
	Name     string
	GuildId  string
	Role     byte
	Accounts []string
}

const memberCName = "member"

const (
	RoleSysAdmin = 1 << iota
	RoleMaster
	RoleAdmin
	RoleNormal
)

func (m *Member) Save() error {
	m.Id = bson.NewObjectId()
	return db.Create(memberCName, m)
}
func (m *Member) UpdateById() error {
	u := bson.M{
		"mobile":  m.Mobile,
		"ability": m.Ability,
		"name":    m.Name,
		"role":    m.Role,
	}
	return db.UpdateOne(memberCName, m.Id, bson.M{"$set": u})
}
func (m *Member) AddAccountById(accId string) error {
	if bson.IsObjectIdHex(accId) {
		u := bson.M{
			"accounts": accId,
		}
		return db.UpdateOne(memberCName, m.Id, bson.M{"$addToSet": u})
	} else {
		return errors.New("invalid account")
	}
}
func FindAllMembers(gId string) ([]Member, error) {
	var ms []Member
	err := db.FindMany(memberCName, bson.M{"guildid": gId}, &ms, "")
	return ms, err
}
func (m *Member) GetById(id bson.ObjectId) error {
	return db.FindOneById(memberCName, id, m)
}
func DelMemberById(id bson.ObjectId) error {
	return db.DeleteOne(memberCName, id)
}
func DelMembersByGId(gId string) error {
	return db.DeleteMany(memberCName, bson.M{"guildid": gId})
}
func roleByAccount(aId, gId string) (byte, error) {
	var ms []Member
	if err := db.FindMany(memberCName, bson.M{"accounts": aId, "guildid": gId}, &ms, ""); err != nil {
		return 0, err
	} else {
		var role byte
		for _, m := range ms {
			role |= m.Role
		}
		return role, nil
	}
}
func RoleAble(sId, gId string) (role byte, err error) {
	role = 0
	err = nil
	var sIdi uint64
	sIdi, err = strconv.ParseUint(sId, 10, 64)
	if err == nil {
		sess := &session.Session{SessId: sIdi}
		if 0 != sIdi {
			err = sess.Get(sIdi)
		}
		if err == nil && bson.IsObjectIdHex(sess.AccountId) {
			//其实为对外的功能权限
			a := Account{}
			err = (&a).GetById(bson.ObjectIdHex(sess.AccountId))
			if err == nil {
				if (&a).IsSysAdmin() {
					role = RoleSysAdmin | RoleAdmin | RoleMaster | RoleNormal
				} else {
					role, err = roleByAccount(a.Id.Hex(), gId)
					if err == nil {
						if 0 != role {
							role |= RoleNormal
						}
						if RoleMaster == role {
							role |= RoleAdmin
						}
					}
				}
			}
		}
	}
	return
}
