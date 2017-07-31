package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/session"
	// "gopkg.in/mgo.v2/bson"
)

//get-account
type GetAccountReq struct {
	SessionId uint64
}

type GetAccountRsp struct {
	AccountId string
	Role      byte
}

func (req *GetAccountReq) GetName() (string, string) {
	return "get-account-req", "get-account-rsp"
}
func (req *GetAccountReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GetAccountRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GetAccountReq) Handle(sess *session.Session) ([]byte, error) {
	rsp := &GetAccountRsp{sess.AccountId, 0}
	// if sess.AccountId.Valid() {
	// 	rsp.AccountId = sess.AccountId.Hex()
	// }
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

//login
type LoginReq struct {
	VerifyCode uint32
}

type LoginRsp struct {
	Result bool
}

func (req *LoginReq) GetName() (string, string) {
	return "login-req", "login-rsp"
}
func (req *LoginReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *LoginRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *LoginReq) Handle(sess *session.Session) ([]byte, error) {
	if req.VerifyCode != sess.VerifyCode {
		return NormalError(ErrVerifyCode)
	}
	ac := &db.Account{}
	if err := ac.GetByMobile(sess.Mobile); err != nil && "not found" != err.Error() {
		return nil, err
	} else {
		if ac.Id == "" {
			ac = &db.Account{Mobile: sess.Mobile}
			if err := ac.Save(); err != nil {
				return nil, err
			}
		}
	}

	if err := sess.SetAccount(ac.Id); err != nil {
		return nil, err
	}
	rsp := &LoginRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

//logout
type LogoutReq struct {
	Mobile string
}

type LogoutRsp struct {
	Result bool
}

func (req *LogoutReq) GetName() (string, string) {
	return "logout-req", "logout-rsp"
}
func (req *LogoutReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *LogoutRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *LogoutReq) Handle(sess *session.Session) ([]byte, error) {
	ac := &db.Account{}
	if err := sess.SetAccount(ac.Id); err != nil {
		return nil, err
	}
	rsp := &LogoutRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
