package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type MemberInfoReq struct {
	Id string
}
type MemberInfoRsp struct {
	Mobile string
}

func (req *MemberInfoReq) GetName() (string, string) {
	return "member-info-req", "member-info-rsp"
}
func (req *MemberInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *MemberInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *MemberInfoReq) Handle(sess *session.Session) ([]byte, error) {
	m := db.Member{}
	if bson.IsObjectIdHex(req.Id) {
		(&m).GetById(bson.ObjectIdHex(req.Id))
	}
	rsp := &MemberInfoRsp{m.Mobile}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

type MemberSaveReq struct {
	Id string
	Mobile string
	GuildId string
}
type MemberSaveRsp struct {
	Result bool
}

func (req *MemberSaveReq) GetName() (string, string) {
	return "member-save-req", "member-save-rsp"
}
func (req *MemberSaveReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *MemberSaveRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *MemberSaveReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		m := &db.Member{Id:bson.ObjectIdHex(req.Id), Mobile:req.Mobile, GuildId:req.GuildId}
		if err := m.UpdateById(); err != nil {
			return nil, err
		}
	} else {
		m := &db.Member{Mobile:req.Mobile, GuildId:req.GuildId}
		if err := m.Save(); err != nil {
			return nil, err
		}
	}
	rsp := &MemberSaveRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

type MemberDeleteReq struct {
	Id string
}
type MemberDeleteRsp struct {
	Result bool
}

func (req *MemberDeleteReq) GetName() (string, string) {
	return "member-delete-req", "member-delete-rsp"
}
func (req *MemberDeleteReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *MemberDeleteRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *MemberDeleteReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		db.DelMemberById(bson.ObjectIdHex(req.Id))
	}
	rsp := &MemberDeleteRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
