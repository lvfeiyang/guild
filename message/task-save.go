package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

type TaskSaveReq struct {
	Id string
	Price int
	Desc string
	GuildId string
}
type TaskSaveRsp struct {
	Result bool
}

func (req *TaskSaveReq) GetName() (string, string) {
	return "task-save-req", "task-save-rsp"
}
func (req *TaskSaveReq) Decode(msgData []byte) error {
	return json.Unmarshal(msgData, req)
}
func (rsp *TaskSaveRsp) Encode() ([]byte, error) {
	return json.Marshal(rsp)
}
func (req *TaskSaveReq) Handle(sess *session.Session) ([]byte, error) {
	if bson.IsObjectIdHex(req.Id) {
		t := &db.Task{Id:bson.ObjectIdHex(req.Id), Price:req.Price, Desc:req.Desc, GuildId:req.GuildId}
		if err := t.UpdateById(); err != nil {
			return nil, err
		}
	} else {
		t := &db.Task{Price:req.Price, Desc:req.Desc, GuildId:req.GuildId}
		if err := t.Save(); err != nil {
			return nil, err
		}
	}
	rsp := &TaskSaveRsp{true}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}
