package message

import (
	"encoding/json"
	// "fmt"
	"github.com/lvfeiyang/guild/common/session"
	// "github.com/lvfeiyang/guild/common/sm"
)

type GetMobileCodeReq struct {
	// SessionId uint64 `json:"-"`
	Mobile string
}

func (req *GetMobileCodeReq) GetName() (string, string) {
	return "get-mobile-code-req", "get-mobile-code-rsp"
}
func (req *GetMobileCodeReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (req *GetMobileCodeReq) Handle(sess *session.Session) ([]byte, error) {
	// nilMsg := []byte{} //无响应 err不能是内部定义的ErrorMsg
	nilMsg := []byte(`{"nil":true}`)

	if err := sess.SetMobile(req.Mobile); err != nil {
		return nil, err
	}
	if 0 == sess.VerifyCode {
		if _, err := sess.SetVerifyCode(); err != nil {
			return nil, err
		}
	}
	// sm.SendVerifyCode(req.Mobile, fmt.Sprintf("%06d", sess.VerifyCode))
	return nilMsg, nil
}
