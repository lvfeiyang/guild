package message

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/session"
	"golang.org/x/crypto/sha3"
)

type ChgPwdReq struct {
	VerifyCode uint32
	Pwd        string
}
type ChgPwdRsp struct {
	Result bool
}

func (req *ChgPwdReq) GetName() (string, string) {
	return "chg-pwd-req", "chg-pwd-rsp"
}
func (req *ChgPwdReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *ChgPwdRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}

//TODO: 接口模块化
func (req *ChgPwdReq) Handle(sess *session.Session) ([]byte, error) {
	ac := &db.Account{}
	if req.VerifyCode != sess.VerifyCode {
		return NormalError(ErrVerifyCode)
	}
	if err := ac.GetByMobile(sess.Mobile); err != nil {
		return nil, err
	}
	if ac.Id == "" {
		return NormalError(ErrNoUser)
	}
	pwdH := sha3.Sum256([]byte(req.Pwd))
	pwdHH := hex.EncodeToString(pwdH[:])
	ac.Pwd = pwdHH
	if err := ac.ChangePwd(); err != nil {
		return nil, err
	}

	rsp := &ChgPwdRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
