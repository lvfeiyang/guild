package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/session"
)

type GetUserInfoReq struct {
}
type GetUserInfoRsp struct {
	Mobile    string
	NickName  string
	Head      string
	Birthday  uint64
	Sex       byte
	Type      byte
	Integrals uint64
}

func (req *GetUserInfoReq) GetName() (string, string) {
	return "get-user-info-req", "get-user-info-rsp"
}
func (req *GetUserInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GetUserInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GetUserInfoReq) Handle(sess *session.Session) ([]byte, error) {
	ac := &db.Account{}
	if err := ac.GetById(sess.AccountId); err != nil {
		return nil, err
	}
	rsp := &GetUserInfoRsp{ac.Mobile, ac.NickName, ac.Head, ac.Birthday, ac.Sex, ac.Type, ac.Integrals}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
