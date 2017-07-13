package db

import "gopkg.in/mgo.v2/bson"

type Task struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Desc     string
	Price    int
	DeadLine int64
	GuildId  string
}

const taskCName = "task"

func (t *Task) Save() error {
	t.Id = bson.NewObjectId()
	return Create(taskCName, t)
}
func (t *Task) UpdateById() error {
	return UpdateOne(taskCName, t.Id, bson.M{"$set": bson.M{"desc": t.Desc, "price": t.Price}})
}
func FindAllTasks(gId string) ([]Task, error) {
	var ts []Task
	err := FindMany(taskCName, bson.M{"guildid": gId}, &ts)
	return ts, err
}
func (t *Task) GetById(id bson.ObjectId) error {
	return FindOneById(taskCName, id, t)
}
func DelTaskById(id bson.ObjectId) error {
	return DeleteOne(taskCName, id)
}
func DelTasksByGId(gId string) error {
	return DeleteMany(taskCName, bson.M{"guildid": gId})
}
