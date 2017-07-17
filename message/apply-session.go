package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
)

type ApplySessionReq struct {
	Device  string
	Ip string
}
type ApplySessionRsp struct {
	SessionId uint64
}

func (req *ApplySessionReq) GetName() (string, string) {
	return "apply-session-req", "apply-session-rsp"
}
func (req *ApplySessionReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *ApplySessionRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (rsp *ApplySessionRsp) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, rsp)
}
func (req *ApplySessionReq) Handle(sess *session.Session) ([]byte, error) {
	client := session.ConnRedis()
	defer client.Close()

	s := &session.Session{}
	sid := s.Apply()
	if 0 == sid {
		return NormalError(ErrNoSession)
	}
	rsp := &ApplySessionRsp{SessionId: sid}
	rspJ, err := rsp.Encode()
	if err != nil {
		return nil, err
	}
	return rspJ, nil
}
