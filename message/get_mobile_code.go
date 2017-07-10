package message

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/sm"
)

type GetMobileCodeReq struct {
	SessionId uint64 `json:"-"`
	Mobile    string
	Sign      string
}

func (req *GetMobileCodeReq) GetName() (string, string) {
	return "get-mobile-code-req", "get-mobile-code-rsp"
}
func (req *GetMobileCodeReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (req *GetMobileCodeReq) SignData() []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 0) //req.SessionId)
	data = append(data, []byte(req.Mobile)...)
	return data
}
func (req *GetMobileCodeReq) Handle(sess *session.Session) ([]byte, error) {
	nilMsg := []byte{} //无响应 err不能是内部定义的ErrorMsg
	/*sess := &session.Session{}
	req.SessionId = sessId
	if err := sess.Get(req.SessionId); err != nil {
		return nil, err
		// return nil, errors.New(CodeMsgMap[ErrGetSessionFail])
	}*/
	recvSig, err := hex.DecodeString(req.Sign)
	if err != nil {
		return nil, err
	}
	if Verify(recvSig, req.SignData(), NewKey(sess.N)) {
		if err := sess.SetMobile(req.Mobile); err != nil {
			return nil, err
		}
		if 0 == sess.VerifyCode {
			if _, err := sess.SetVerifyCode(); err != nil {
				return nil, err
			}
		}
		sm.SendVerifyCode(req.Mobile, fmt.Sprintf("%06d", sess.VerifyCode))
	} else {
		return nil, errors.New(CodeMsgMap[ErrNoVerify])
	}
	return nilMsg, nil
}
