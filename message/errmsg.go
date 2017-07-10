package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/flog"
)

const (
	ErrNoVerify       = 1
	ErrDecode         = 2
	ErrUnknow         = 3
	ErrNoSession      = 4
	ErrGetSessionFail = 5
	ErrUnknowMsg      = 6
	ErrUsedClientN    = 7
	ErrDeCrypto       = 8
	ErrVerifyCode     = 9
	ErrPwdWrong       = 10
	ErrHavenRegister  = 11
	ErrNoUser         = 12
)

var CodeMsgMap = map[uint32]string{
	ErrNoVerify:       "sign can not verify!",
	ErrUnknow:         "unknow error happen",
	ErrNoSession:      "no session",
	ErrGetSessionFail: "get session fail",
	ErrUnknowMsg:      "recv unknow msg",
	ErrUsedClientN:    "client n have used",
	ErrDeCrypto:       "decode crypto fail",
	ErrVerifyCode:     "invalid verify code",
	ErrPwdWrong:       "get a wrong password",
	ErrHavenRegister:  "mobile have been used",
	ErrNoUser:         "have not this user",
}

type ErrorMsg struct {
	ErrCode uint32
	ErrMsg  string
}

func (errMsg *ErrorMsg) Encode() ([]byte, error) {
	return json.Marshal(errMsg)
}

func (errMsg *ErrorMsg) Decode(errData []byte) error {
	return json.Unmarshal(errData, errMsg)
}

func (errMsg *ErrorMsg) Error() string {
	return errMsg.ErrMsg
}

func NormalError(errCode uint32) ([]byte, error) {
	errMsg := &ErrorMsg{ErrCode: errCode, ErrMsg: CodeMsgMap[errCode]}
	if errData, err := errMsg.Encode(); err != nil {
		// flog.LogFile.Println(err)
		return nil, err
	} else {
		return errData, errMsg
	}
}

func DecodeError(msgName string) string {
	errMsg := &ErrorMsg{ErrCode: ErrDecode, ErrMsg: "decode " + msgName + " fail"}
	if errData, err := errMsg.Encode(); err != nil {
		flog.LogFile.Println(err)
		return ""
	} else {
		return string(errData)
	}
}

func UnknowError() string {
	errMsg := &ErrorMsg{ErrCode: ErrUnknow, ErrMsg: CodeMsgMap[ErrUnknow]}
	if errData, err := errMsg.Encode(); err != nil {
		flog.LogFile.Println(err)
		return ""
	} else {
		return string(errData)
	}
}

func UnknowMsg() string {
	errMsg := &ErrorMsg{ErrCode: ErrUnknowMsg, ErrMsg: CodeMsgMap[ErrUnknowMsg]}
	if errData, err := errMsg.Encode(); err != nil {
		flog.LogFile.Println(err)
		return ""
	} else {
		return string(errData)
	}
}
