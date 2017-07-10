package message

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
)

type ChgUserInfoReq struct {
	NickName string
	Head     string
	Birthday uint64
	Sex      uint8
	Type     uint8
	Sign     string
}
type ChgUserInfoRsp struct {
	Result bool
}

func (req *ChgUserInfoReq) GetName() (string, string) {
	return "chg-user-info-req", "chg-user-info-rsp"
}
func (req *ChgUserInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (req *ChgUserInfoReq) SignData() []byte {
	data := make([]byte, 10)
	binary.LittleEndian.PutUint64(data[:8], req.Birthday)
	data[8] = req.Sex
	data[9] = req.Type
	data = append(data, []byte(req.NickName+req.Head)...)
	return data
}
func (rsp *ChgUserInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *ChgUserInfoReq) Handle(sess *session.Session) ([]byte, error) {
	nkey := NewKey(sess.N)
	recvSig, err := hex.DecodeString(req.Sign)
	if err != nil {
		return nil, err
	}
	if Verify(recvSig, req.SignData(), nkey) {

	} else {
		return NormalError(ErrNoVerify)
	}
	return nil, nil
}
