package message

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lvfeiyang/guild/common/flog"
	"github.com/lvfeiyang/guild/common/session"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Message struct {
	Name      string
	Data      string
	SessionId uint64
}

func (msg *Message) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg.Name = r.URL.Path[1:] + "-req"//strings.TrimLeft(r.URL.Path, "/") + "-req"
	if cookie, err := r.Cookie(session.CookieKey); err != nil {
		flog.LogFile.Println(err)
	} else {
		if msg.SessionId, err = strconv.ParseUint(cookie.Value, 16, 64); err != nil {
			flog.LogFile.Println(err)
		}
	}
	if 0 == strings.Compare("application/json", r.Header.Get("Content-Type")) {
		defer r.Body.Close()
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			flog.LogFile.Println(err)
		}
		msg.Data = string(buff)

		sendMsg := msg.HandleMsg()
		w.Header().Set("Content-Type", "application/json")
		if 0 == strings.Compare("error-msg", sendMsg.Name) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			if 0 == strings.Compare("apply-session-rsp", sendMsg.Name) {
				asRsp := &ApplySessionRsp{}
				if err := asRsp.Decode([]byte(sendMsg.Data)); err == nil && asRsp.SessionId != 0 {
					cookie := http.Cookie{Name: session.CookieKey, Value: strconv.FormatUint(asRsp.SessionId, 16)}
					http.SetCookie(w, &cookie)
				}
			}
		}
		w.Write(httpFormat(sendMsg))
		// w.Write([]byte(sendMsg.Data))
	} else {
		// IDEA: form表单需整合为json
		return
	}
}
func httpFormat(msg *Message) []byte {
	format := struct {
		Code    uint32
		Message string
		Data    string
	}{}
	if 0 == strings.Compare("error-msg", msg.Name) {
		err := &ErrorMsg{}
		err.Decode([]byte(msg.Data))
		format.Code = err.ErrCode
		format.Message = err.ErrMsg
	} else {
		format.Data = msg.Data
	}
	sendData, _ := json.Marshal(format)
	return sendData
}

func (msg *Message) Decode(data []byte) error {
	return json.Unmarshal(data, msg)
}
func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}

type msgHandleIF interface {
	Decode(msgData []byte) error
	Handle(sess *session.Session) ([]byte, error)
	GetName() (string, string)
}

func deCrypto(msgData []byte, sess *session.Session) ([]byte, error) {
	/* sess = &session.Session{}
	if err := sess.Get(sessId); err != nil {
		// flog.LogFile.Println(err)
		return nil, err
	}*/
	recvEn := make([]byte, hex.DecodedLen(len(msgData)))
	n, err := hex.Decode(recvEn, msgData)
	if err != nil {
		return nil, err
	}
	recv, err := AesDe(recvEn[:n], NewKey(sess.N))
	if err != nil {
		return nil, err
	}
	return recv, nil
}
func handleOneMsg(req msgHandleIF, msgData []byte, sess *session.Session) *Message {
	sendMsg := &Message{Name: "error-msg", Data: UnknowError()}
	reqName, rspName := req.GetName()

	if req.Decode(msgData) != nil {
		sendMsg = &Message{Name: "error-msg", Data: DecodeError(reqName)}
	} else {
		var rspData []byte
		var err interface{}
		// req.SessionId = msgSessId
		rspData, err = req.Handle(sess)
		if err != nil {
			if _, ok := err.(*ErrorMsg); ok {
				sendMsg = &Message{Name: "error-msg", Data: string(rspData)}
			} else {
				flog.LogFile.Println(err)
			}
		} else {
			sendMsg = &Message{Name: rspName, Data: string(rspData)}
		}
	}
	return sendMsg
}

// IDEA: 可改为 rpc 做转发分担 不断线升级
func (msg *Message) HandleMsg() *Message {
	// IDEA: 用接口定义去掉 switch case
	sess := &session.Session{SessId: msg.SessionId}
	if 0 != msg.SessionId {
		if err := sess.Get(msg.SessionId); err != nil {
			errData, _ := NormalError(ErrGetSessionFail)
			return &Message{Name: "error-msg", Data: string(errData)}
		}
	}
	switch msg.Name {
	case "apply-session-req":
		return handleOneMsg(&ApplySessionReq{}, []byte(msg.Data), sess)
	case "get-n-req":
		return handleOneMsg(&GetNReq{}, []byte(msg.Data), sess)
	case "get-mobile-code-req":
		return handleOneMsg(&GetMobileCodeReq{}, []byte(msg.Data), sess)
	case "register-req":
		msgData, err := deCrypto([]byte(msg.Data), sess)
		if err != nil {
			errData, _ := NormalError(ErrDeCrypto)
			return &Message{Name: "error-msg", Data: string(errData)}
		}
		return handleOneMsg(&RegisterReq{}, msgData, sess)
	case "login-req":
		msgData, err := deCrypto([]byte(msg.Data), sess)
		if err != nil {
			errData, _ := NormalError(ErrDeCrypto)
			return &Message{Name: "error-msg", Data: string(errData)}
		}
		return handleOneMsg(&LoginReq{}, msgData, sess)
	default:
		return &Message{Name: "error-msg", Data: UnknowMsg()}
	}
}
