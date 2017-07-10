package message

import (
	"encoding/binary"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
)

type GetNReq struct {
	SessionId uint64
}
type GetNRsp struct {
	RandomN uint64
	Sign    string
}

func (req *GetNReq) GetName() (string, string) {
	return "get-n-req", "get-n-rsp"
}
func (req *GetNReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *GetNRsp) SignData() []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, rsp.RandomN)
	return data
}
func (rsp *GetNRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *GetNReq) Handle(sess *session.Session) ([]byte, error) {
	/*sess := &session.Session{}
	if err := sess.Get(sessId); err != nil {
		return NormalError(ErrGetSessionFail)
	}*/
	rsp := &GetNRsp{RandomN: sess.N}
	nKey := NewKey(sess.N)

	sign, err := Sign(rsp.SignData(), nKey)
	if err != nil {
		return nil, err
	}
	rsp.Sign = string(sign)
	rspJ, err := rsp.Encode()
	if err != nil {
		return nil, err
	}
	return AesEn(rspJ, nKey)
}
