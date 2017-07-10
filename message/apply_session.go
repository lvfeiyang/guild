package message

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"strconv"
	"time"
)

type ApplySessionReq struct {
	Device  string
	ClientN uint64
	Sign    string
}
type ApplySessionRsp struct {
	SessionId uint64
	RandomN   uint64
	Sign      string
}

func (req *ApplySessionReq) GetName() (string, string) {
	return "apply-session-req", "apply-session-rsp"
}
func (req *ApplySessionReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (req *ApplySessionReq) SignData() []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data[:8], req.ClientN)
	data = append(data, []byte(req.Device)...)
	return data
}
func (rsp *ApplySessionRsp) SignData() []byte {
	data := make([]byte, 16)
	binary.LittleEndian.PutUint64(data[:8], rsp.SessionId)
	binary.LittleEndian.PutUint64(data[8:], rsp.RandomN)
	return data
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
	if res, err := client.SetNX(session.PrefixClientN+strconv.FormatUint(req.ClientN, 10), 1, 8*time.Hour).Result(); err != nil {
		return nil, err
	} else {
		if false && false == res {
			return NormalError(ErrUsedClientN)
		}
	}

	clientKey := NewKey(req.ClientN)
	recvSig, err := hex.DecodeString(req.Sign)
	if err != nil {
		return nil, err
	}
	if true || Verify(recvSig, req.SignData(), clientKey) {
		// var s *session.Session
		s := &session.Session{}
		sid := s.Apply()
		if 0 == sid {
			return NormalError(ErrNoSession)
		}
		rsp := &ApplySessionRsp{SessionId: sid, RandomN: s.N}

		sign, err := Sign(rsp.SignData(), NewKey(s.N))
		if err != nil {
			return nil, err
		}
		// rsp.Sign = string(sign)
		rsp.Sign = hex.EncodeToString(sign)
		rspJ, err := rsp.Encode()
		if err != nil {
			return nil, err
		}
		return rspJ, nil
		enData, err := AesEn(rspJ, clientKey)
		if err != nil {
			return nil, err
		}
		send := make([]byte, hex.EncodedLen(len(enData)))
		hex.Encode(send, enData)
		return send, nil
	} else {
		return NormalError(ErrNoVerify)
	}
	// return nil, nil //[]byte{}
}
