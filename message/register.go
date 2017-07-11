package message

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/session"
	"golang.org/x/crypto/sha3"
)

type RegisterReq struct {
	VerifyCode uint32
	Pwd        string
	NickName   string
	Head       string
}

type RegisterRsp struct {
	Result bool
}

func (req *RegisterReq) GetName() (string, string) {
	return "register-req", "register-rsp"
}
func (req *RegisterReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *RegisterRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *RegisterReq) Handle(sess *session.Session) ([]byte, error) {
	//TODO: 绑定 C+e 到错误代码块粘贴
	if req.VerifyCode != sess.VerifyCode {
		return NormalError(ErrVerifyCode)
	}
	ac := &db.Account{}
	if err := ac.GetByMobile(sess.Mobile); err != nil {
		return nil, err
	}
	if ac.Id != "" {
		return NormalError(ErrHavenRegister)
	}
	pwdH := sha3.Sum256([]byte(req.Pwd))
	pwdHH := hex.EncodeToString(pwdH[:])
	ac = &db.Account{Pwd: pwdHH, Mobile: sess.Mobile}
	if err := ac.Save(); err != nil {
		return nil, err
	}
	// if err := sess.SetStatus(session.StatusLogin); err != nil {
	if err := sess.SetAccount(ac.Id); err != nil {
		return nil, err
	}
	rsp := &RegisterRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
