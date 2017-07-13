package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type TaskDeleteReq struct {
	Id string
}
type TaskDeleteRsp struct {
	Result bool
}

func (req *TaskDeleteReq) GetName() (string, string) {
	return "task-delete-req", "task-delete-rsp"
}
func (req *TaskDeleteReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *TaskDeleteRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *TaskDeleteReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		db.DelTaskById(bson.ObjectIdHex(req.Id))
	}
	rsp := &TaskDeleteRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
