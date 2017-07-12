package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type GuildDeleteReq struct {
	Id string
}
type GuildDeleteRsp struct {
	Result bool
}

func (req *GuildDeleteReq) GetName() (string, string) {
	return "guild-delete-req", "guild-delete-rsp"
}
func (req *GuildDeleteReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GuildDeleteRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GuildDeleteReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		db.DelGuildById(bson.ObjectIdHex(req.Id))
	}
	rsp := &GuildDeleteRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
