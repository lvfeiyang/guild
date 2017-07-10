package message

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/session"
	"golang.org/x/crypto/sha3"
	"strings"
)

type LoginReq struct {
	Mobile string
	Pwd    string
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
	ac := &db.Account{}
	if err := ac.GetByMobile(req.Mobile); err != nil {
		return nil, err
	}
	if ac.Id == "" {
		return NormalError(ErrHavenRegister)
	}
	pwdH := sha3.Sum256([]byte(req.Pwd))
	pwdHH := hex.EncodeToString(pwdH[:])
	if 0 != strings.Compare(pwdHH, ac.Pwd) {
		return NormalError(ErrPwdWrong)
	}
	// if err := sess.SetStatus(session.StatusLogin); err != nil {
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
