package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/proxy/common/session"
	"github.com/lvfeiyang/proxy/message"
	// "gopkg.in/mgo.v2/bson"
	"strconv"
)

//get-account
type GetAccountReq struct {
	SessionId uint64
	GuildId   string
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
	/*if bson.IsObjectIdHex(sess.AccountId) {
		//其实为对外的功能权限
		a := db.Account{}
		if err := (&a).GetById(bson.ObjectIdHex(sess.AccountId)); err != nil {
			return nil, err
		}
		if (&a).IsSysAdmin() {
			rsp.Role = db.RoleSysAdmin | db.RoleAdmin | db.RoleMaster
		} else {
			if r, err := db.RoleByAccount(a.Id.Hex(), req.GuildId); err != nil {
				return nil, err //TODO 考虑使用panic recover defer
			} else {
				if db.RoleMaster == r {
					rsp.Role = db.RoleMaster | db.RoleAdmin
				} else {
					rsp.Role = r
				}
			}
		}
	}*/
	rsp.Role, _ = db.RoleAble(strconv.FormatUint(sess.SessId, 10), req.GuildId)
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
		return message.NormalError(message.ErrVerifyCode)
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
