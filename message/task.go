package message

import (
	"encoding/json"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/proxy/common/session"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type TaskInfoReq struct {
	Id string
}
type TaskInfoRsp struct {
	Price    string
	Desc     string
	DeadLine int64
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
	rsp := &TaskInfoRsp{strconv.Itoa(t.Price), t.Desc, t.DeadLine}
	if rspJ, err := rsp.Encode(); err != nil {
		return nil, err
	} else {
		return rspJ, nil
	}
}

type TaskSaveReq struct {
	Id       string
	Price    int
	Desc     string
	DeadLine int64
	GuildId  string
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
		t := &db.Task{Id: bson.ObjectIdHex(req.Id), Price: req.Price, Desc: req.Desc, GuildId: req.GuildId, DeadLine: req.DeadLine}
		if err := t.UpdateById(); err != nil {
			return nil, err
		}
	} else {
		t := &db.Task{Price: req.Price, Desc: req.Desc, GuildId: req.GuildId, DeadLine: req.DeadLine}
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
