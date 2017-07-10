package main

import (
	"net/http"
	"html/template"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/flog"
	"github.com/lvfeiyang/guild/common/config"
	"github.com/lvfeiyang/guild/message"
	"gopkg.in/mgo.v2/bson"
	"path/filepath"
)

var htmlPath string

func main() {
	flog.Init()
	config.Init()
	// session.Init()
	db.Init()
	htmlPath = config.ConfigVal.HtmlPath // E:\leonshare\go-workspace\src\github.com\lvfeiyang

	// var jsFiles, cssFiles string
	// if "linux" == runtime.GOOS {
	// 	jsFiles, cssFiles = htmlPath+"sfk/js", htmlPath+"sfk/css"
	// } else {
	// 	jsFiles, cssFiles = htmlPath+"sfk\\js", htmlPath+"sfk\\css"
	// }
	jsFiles, cssFiles := filepath.Join(htmlPath, "sfk", "js"), filepath.Join(htmlPath, "sfk", "css")
	gcssFiles := filepath.Join(htmlPath, "guild", "html", "css")
	gjsFiles := filepath.Join(htmlPath, "guild", "html", "js")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
	http.Handle("/guild-css/", http.StripPrefix("/guild-css/", http.FileServer(http.Dir(gcssFiles))))
	http.Handle("/guild-js/", http.StripPrefix("/guild-js/", http.FileServer(http.Dir(gjsFiles))))

	http.HandleFunc("/guild", guildHandler)
	http.HandleFunc("/msg/", &message.Message{})//guild-save

	// http.HandleFunc("/guild/edit/", guildEditHandler)
	// http.HandleFunc("/guild/save/", guildSaveHandler)
	// http.HandleFunc("/guild/list", guildListHandler)
	//
	// http.HandleFunc("/member/edit/", memberEditHandler)
	// http.HandleFunc("/member/save/", memberSaveHandler)
	// http.HandleFunc("/member/list/", memberListHandler)
	//
	// http.HandleFunc("/task/edit/", taskEditHandler)
	// http.HandleFunc("/task/save/", taskSaveHandler)
	// http.HandleFunc("/task/list/", taskListHandler)
	//
	// http.HandleFunc("/apply/edit/", applyEditHandler)
	// http.HandleFunc("/apply/save/", applySaveHandler)
	// http.HandleFunc("/apply/list/", applyListHandler)

	flog.LogFile.Fatal(http.ListenAndServe(":80", nil))
}
func guildHandler(w http.ResponseWriter, r *http.Request)  {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "guild.html"),
		filepath.Join(htmlPath, "guild", "html", "sidebar.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "main.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "edit-guild.tmpl"),
	}
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		if err := t.Execute(w, nil); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
