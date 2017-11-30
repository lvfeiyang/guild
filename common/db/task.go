package db

import "gopkg.in/mgo.v2/bson"
import "github.com/lvfeiyang/proxy/common/db"

type Task struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Desc      string
	Price     int
	DeadLine  int64
	GuildId   string
	Client    string
	Recipient []string
}

const taskCName = "task"

func (t *Task) Save() error {
	t.Id = bson.NewObjectId()
	return db.Create(taskCName, t)
}
func (t *Task) UpdateById() error {
	ud := bson.M{"desc": t.Desc, "price": t.Price, "deadline": t.DeadLine}
	return db.UpdateOne(taskCName, t.Id, bson.M{"$set": ud})
}
func FindAllTasks(gId string) ([]Task, error) {
	var ts []Task
	err := db.FindMany(taskCName, bson.M{"guildid": gId}, &ts, db.Option{})
	return ts, err
}
func (t *Task) GetById(id bson.ObjectId) error {
	return db.FindOneById(taskCName, id, t)
}
func DelTaskById(id bson.ObjectId) error {
	return db.DeleteOne(taskCName, id)
}
func DelTasksByGId(gId string) error {
	return db.DeleteMany(taskCName, bson.M{"guildid": gId})
}
