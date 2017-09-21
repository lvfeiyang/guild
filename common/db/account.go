package db

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
	"github.com/lvfeiyang/proxy/common/db"
)

type Account struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Mobile string
}

const accountCName = "account"

func (ac *Account) Save() error {
	ac.Id = bson.NewObjectId()
	return db.Create(accountCName, ac)
}
func (ac *Account) GetByMobile(mobile string) error {
	return db.FindOne(accountCName, bson.M{"mobile": mobile}, ac)
}
func (ac *Account) GetById(id bson.ObjectId) error {
	return db.FindOneById(accountCName, id, ac)
}
func (ac *Account) IsSysAdmin() bool {
	return strings.EqualFold("13917287994", ac.Mobile)
}
