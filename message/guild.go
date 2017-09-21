package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/proxy/common/session"
	"gopkg.in/mgo.v2/bson"
)

type GuildInfoReq struct {
	Id string
}
type GuildInfoRsp struct {
	Name      string
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

type GuildSaveReq struct {
	Id        string
	Name      string
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
		g := &db.Guild{Id: bson.ObjectIdHex(req.Id), Name: req.Name, Introduce: req.Introduce}
		if err := g.UpdateById(); err != nil {
			// flog.LogFile.Println(err)
			return nil, err
		}
	} else {
		g := &db.Guild{Name: req.Name, Introduce: req.Introduce}
		if err := g.Save(); err != nil {
			return nil, err
		}
	}
	rsp := &GuildSaveRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

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
		db.DelTasksByGId(req.Id)
		db.DelMembersByGId(req.Id)
	}
	rsp := &GuildDeleteRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
