package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/proxy/common/session"
	"gopkg.in/mgo.v2/bson"
)

type MemberInfoReq struct {
	Id string
}
type MemberInfoRsp struct {
	Name    string
	Mobile  string
	Ability string
	Role    byte
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
	rsp := &MemberInfoRsp{m.Name, m.Mobile, m.Ability, m.Role}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

type MemberSaveReq struct {
	Id      string
	GuildId string
	Name    string
	Mobile  string
	Ability string
	Role    byte
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
	var m *db.Member
	if bson.IsObjectIdHex(req.Id) {
		m = &db.Member{Id: bson.ObjectIdHex(req.Id), Mobile: req.Mobile, GuildId: req.GuildId, Name: req.Name, Ability: req.Ability, Role: req.Role}
		if err := m.UpdateById(); err != nil {
			return nil, err
		}
	} else {
		m = &db.Member{Mobile: req.Mobile, GuildId: req.GuildId, Name: req.Name, Ability: req.Ability, Role: req.Role}
		if err := m.Save(); err != nil {
			return nil, err
		}
	}
	ac := &db.Account{}
	if err := ac.GetByMobile(req.Mobile); err != nil && "not found" != err.Error() {
		return nil, err
	}
	if err := m.AddAccountById(ac.Id.Hex()); err != nil {
		return nil, err
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
