package main

import (
	"net/http"
	"html/template"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/flog"
	"github.com/lvfeiyang/guild/common/config"
	"github.com/lvfeiyang/guild/message"
	"path/filepath"
	"gopkg.in/mgo.v2/bson"
)

var htmlPath string

func main() {
	flog.Init()
	config.Init()
	// session.Init()
	db.Init()
	htmlPath = config.ConfigVal.HtmlPath // E:\leonshare\go-workspace\src\github.com\lvfeiyang

	jsFiles, cssFiles := filepath.Join(htmlPath, "sfk", "js"), filepath.Join(htmlPath, "sfk", "css")
	gcssFiles := filepath.Join(htmlPath, "guild", "html", "css")
	gjsFiles := filepath.Join(htmlPath, "guild", "html", "js")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
	http.Handle("/guild-css/", http.StripPrefix("/guild-css/", http.FileServer(http.Dir(gcssFiles))))
	http.Handle("/guild-js/", http.StripPrefix("/guild-js/", http.FileServer(http.Dir(gjsFiles))))

	http.Handle("/msg/", &message.Message{})//guild-save

	http.HandleFunc("/guild", guildHandler)
	http.HandleFunc("/guild/detail", guildDetailHandler)
	http.HandleFunc("/task", taskHandler)

	flog.LogFile.Fatal(http.ListenAndServe(":80", nil))
}
func guildHandler(w http.ResponseWriter, r *http.Request)  {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "guild.html"),
		filepath.Join(htmlPath, "guild", "html", "sidebar.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "main.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "table.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "edit-guild.tmpl"),
	}
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		var view struct {
			GuildList []struct {
				Id string
				Name string
			}
		}
		gs, err := db.FindAllGuilds()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range gs {
			view.GuildList = append(view.GuildList, struct{Id, Name string}{v.Id.Hex(), v.Name})
		}
		//t.ExecuteTemplate(w, "sidebar", struct{GuildList []oneViewGuild}{viewGuildList})
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func guildDetailHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "main.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "table.tmpl"),
	}
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		id := r.Form.Get("Id")
		g := db.Guild{}
		if bson.IsObjectIdHex(id) {
			(&g).GetById(bson.ObjectIdHex(id))
		}
		view := struct {
			Id string
			Name string
			Introduce string
		}{g.Id.Hex(), g.Name, g.Introduce}
		if err := t.ExecuteTemplate(w, "main", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func taskHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.ParseFiles(filepath.Join(htmlPath, "guild", "html", "table.tmpl")); err != nil {
		flog.LogFile.Println(err)
	} else {
		view := struct {
			Thead []string
			Tbody []struct {
				Id string
				Desc string
				Price string
			}
		} {Thead: []string{"编号", "描述", "价格"}}
		if err := t.ExecuteTemplate(w, "table", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
