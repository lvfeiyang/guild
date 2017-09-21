package message

import (
	"github.com/lvfeiyang/proxy/message"
	"net/http"
)

var MhMap map[string]message.MsgHandleIF

func Init() {
	MhMap = map[string]message.MsgHandleIF{
		"guild-save-req":      &GuildSaveReq{},
		"guild-info-req":      &GuildInfoReq{},
		"guild-delete-req":    &GuildDeleteReq{},
		"task-save-req":       &TaskSaveReq{},
		"task-info-req":       &TaskInfoReq{},
		"task-delete-req":     &TaskDeleteReq{},
		"member-save-req":     &MemberSaveReq{},
		"member-info-req":     &MemberInfoReq{},
		"member-delete-req":   &MemberDeleteReq{},
		"apply-session-req":   &ApplySessionReq{},
		"get-n-req":           &GetNReq{},
		"get-mobile-code-req": &GetMobileCodeReq{},
		"login-req":           &LoginReq{},
		"logout-req":          &LogoutReq{},
		"get-account-req":     &GetAccountReq{},
	}
	return
}

type LocMessage message.Message

func (msg *LocMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message.GeneralServeHTTP((*message.Message)(msg), w, r, MhMap)
}
