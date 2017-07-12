package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type GuildInfoReq struct {
	Id string
}
type GuildInfoRsp struct {
	Name string
	Introduce string
}

func (req *GuildInfoReq) GetName() (string, string) {
	return "guild-info-req", "guild-info-rsp"
}
func (req *GuildInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GuildInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GuildInfoReq) Handle(sess *session.Session) ([]byte, error) {
	g := db.Guild{}
	if bson.IsObjectIdHex(req.Id) {
		(&g).GetById(bson.ObjectIdHex(req.Id))
	}
	rsp := &GuildInfoRsp{g.Name, g.Introduce}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
