package main

import (
	"log"
	"net/http"
	"html/template"
	"github.com/lvfeiyang/guild/common/db"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	http.HandleFunc("/guild/edit", guildEditHandler)

	log.Fatal(http.ListenAndServe(":7777", nil))
}
func guildEditHandler(w http.ResponseWriter, r *http.Request)  {
	if t, err := template.ParseFiles("html/guild-edit.html"); err != nil {
		log.Println(err)
	} else {
		id := r.URL.Path[len("/guild/edit/"):]
		g := db.Guild{}
		(&g).GetById(bson.ObjectIdHex(id))
		if err := t.Execute(w, g); err != nil {
			log.Println(err)
		}
	}
}
