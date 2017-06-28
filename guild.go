package main

import (
	// "log"
	"net/http"
	"html/template"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/flog"
	"github.com/lvfeiyang/guild/common/config"
	"gopkg.in/mgo.v2/bson"
)

var htmlPath = "C:\\Users\\lxm19\\workspace\\go\\src\\github.com\\lvfeiyang\\guild\\"

func main() {
	flog.Init()
	config.Init()
	// session.Init()
	db.Init()

	http.HandleFunc("/guild/edit/", guildEditHandler)
	http.HandleFunc("/guild/save/", guildSaveHandler)
	http.HandleFunc("/guild/list", guildListHandler)

	http.HandleFunc("/member/edit/", memberEditHandler)
	http.HandleFunc("/member/save/", memberSaveHandler)
	http.HandleFunc("/member/list/", memberListHandler)

	http.HandleFunc("/task/edit/", taskEditHandler)
	http.HandleFunc("/task/save/", taskSaveHandler)
	http.HandleFunc("/task/list/", taskListHandler)

	http.HandleFunc("/apply/edit/", applyEditHandler)
	http.HandleFunc("/apply/save/", applySaveHandler)
	http.HandleFunc("/apply/list/", applyListHandler)

	flog.LogFile.Fatal(http.ListenAndServe(":80", nil))
}
func guildEditHandler(w http.ResponseWriter, r *http.Request)  {
	if t, err := template.ParseFiles(htmlPath+"html\\guild-edit.html"); err != nil {
		flog.LogFile.Println(err)
	} else {
		id := r.URL.Path[len("/guild/edit/"):]
		g := db.Guild{}
		if bson.IsObjectIdHex(id) {
			(&g).GetById(bson.ObjectIdHex(id))
		}
		view := struct {
			Id string
			Name string
			Introduce string
		}{g.Id.Hex(), g.Name, g.Introduce}
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func guildSaveHandler(w http.ResponseWriter, r *http.Request)  {
	id := r.URL.Path[len("/guild/edit/"):]
	name := r.FormValue("name")
	introduce := r.FormValue("introduce")
	if bson.IsObjectIdHex(id) {
		g := &db.Guild{Id:bson.ObjectIdHex(id), Name:name, Introduce:introduce}
		if err := g.UpdateById(); err != nil {
			flog.LogFile.Println(err)
		}
	} else {
		g := &db.Guild{Name:name, Introduce:introduce}
		if err := g.Save(); err != nil {
			flog.LogFile.Println(err)
		}
	}
	http.Redirect(w, r, "/guild/list", http.StatusFound)
}
func guildListHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.ParseFiles(htmlPath+"html\\guild-list.html"); err != nil {
		flog.LogFile.Println(err)
	} else {
		type oneView struct {
			Id string
			Name string
			Introduce string
		}
		var view []oneView
		gs, err := db.FindAllGuilds()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range gs {
			view = append(view, oneView{v.Id.Hex(), v.Name, v.Introduce})
		}
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
		}
	}
}
