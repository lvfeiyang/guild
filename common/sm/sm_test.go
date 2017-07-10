package sm

import (
	"github.com/lvfeiyang/guild/common/flog"
	"testing"
	"time"
)

func TestSendVerifyCode(t *testing.T) {
	flog.Init()
	SendVerifyCode("13917287994", "123456")
	time.Sleep(5 * time.Second)
	// sendSM("13917287994", "1", []string{"123456"})
}
