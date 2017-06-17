package config

import (
	"encoding/json"
	"github.com/lvfeiyang/ss/common/flog"
	"io/ioutil"
	"runtime"
)

type config struct {
	ConnectType string
	RedisUrl    string
	MongoUrl    string
}

var ConfigVal = &config{}

func Init() {
	var filePath string
	if "linux" == runtime.GOOS {
		filePath = "/root/guild/config"
	} else {
		filePath = "C:\\Users\\lxm19\\config"
	}
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		flog.LogFile.Fatal(err)
	}
	err = json.Unmarshal(conf, ConfigVal)
	if err != nil {
		flog.LogFile.Fatal(err)
	}
}