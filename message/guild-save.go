package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type GuildSaveReq struct {
	Id string
	Name string
	Introduce string
}
type GuildSaveRsp struct {
	Result bool
}

func (req *GuildSaveReq) GetName() (string, string) {
	return "guild-save-req", "guild-save-rsp"
}
func (req *GuildSaveReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GuildSaveRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GuildSaveReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		g := &db.Guild{Id:bson.ObjectIdHex(req.Id), Name:req.Name, Introduce:req.Introduce}
		if err := g.UpdateById(); err != nil {
			// flog.LogFile.Println(err)
			return nil, err
		}
	} else {
		g := &db.Guild{Name:req.Name, Introduce:req.Introduce}
		if err := g.Save(); err != nil {
			return nil, err
		}
	}
	rsp := &LoginRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
