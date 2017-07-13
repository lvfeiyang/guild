package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type TaskInfoReq struct {
	Id string
}
type TaskInfoRsp struct {
	Price string
	Desc string
}

func (req *TaskInfoReq) GetName() (string, string) {
	return "task-info-req", "task-info-rsp"
}
func (req *TaskInfoReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *TaskInfoRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *TaskInfoReq) Handle(sess *session.Session) ([]byte, error) {
	t := db.Task{}
	if bson.IsObjectIdHex(req.Id) {
		(&t).GetById(bson.ObjectIdHex(req.Id))
	}
	rsp := &TaskInfoRsp{strconv.Itoa(t.Price), t.Desc}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
