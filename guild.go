package main

import (
	"github.com/lvfeiyang/guild/common/config"
	"github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/guild/common/flog"
	"github.com/lvfeiyang/guild/common/session"
	"github.com/lvfeiyang/guild/message"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

var htmlPath string

func main() {
	flog.Init()
	config.Init()
	session.Init()
	db.Init()
	htmlPath = config.ConfigVal.HtmlPath // E:\leonshare\go-workspace\src\github.com\lvfeiyang

	jsFiles, cssFiles, fontsFiles := filepath.Join(htmlPath, "sfk", "js"), filepath.Join(htmlPath, "sfk", "css"), filepath.Join(htmlPath, "sfk", "fonts")
	gcssFiles := filepath.Join(htmlPath, "guild", "html", "css")
	gjsFiles := filepath.Join(htmlPath, "guild", "html", "js")
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(fontsFiles))))
	http.Handle("/guild-css/", http.StripPrefix("/guild-css/", http.FileServer(http.Dir(gcssFiles))))
	http.Handle("/guild-js/", http.StripPrefix("/guild-js/", http.FileServer(http.Dir(gjsFiles))))

	http.Handle("/msg/", &message.Message{}) //guild-save

	http.HandleFunc("/guild", guildHandler)
	http.HandleFunc("/guild/detail", guildDetailHandler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/member", memberHandler)

	flog.LogFile.Fatal(http.ListenAndServe(":80", nil))
}
func guildHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "guild.html"),
		filepath.Join(htmlPath, "guild", "html", "sidebar.tmpl"),
		// filepath.Join(htmlPath, "guild", "html", "main.tmpl"),
		// filepath.Join(htmlPath, "guild", "html", "task-table.tmpl"),
		// filepath.Join(htmlPath, "guild", "html", "member-table.tmpl"),

		filepath.Join(htmlPath, "guild", "html", "modal", "edit-guild.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "edit-member.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "edit-task.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "login.tmpl"),
	}
	// pattern := filepath.Join(htmlPath, "guild", "html", "modal", "*.tmpl");
	// t, err := template.ParseGlob(pattern)
	// if err != nil {
	// 	flog.LogFile.Println(err)
	// }
	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		var view struct {
			GuildList []struct {
				Id   string
				Name string
			}
		}
		gs, err := db.FindAllGuilds()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range gs {
			view.GuildList = append(view.GuildList, struct{ Id, Name string }{v.Id.Hex(), v.Name})
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
		// filepath.Join(htmlPath, "guild", "html", "task-table.tmpl"),
		// filepath.Join(htmlPath, "guild", "html", "member-table.tmpl"),
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
			Id        string
			Name      string
			Introduce string
		}{g.Id.Hex(), g.Name, g.Introduce}
		if err := t.ExecuteTemplate(w, "main", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func taskHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.ParseFiles(filepath.Join(htmlPath, "guild", "html", "task-table.tmpl")); err != nil {
		flog.LogFile.Println(err)
	} else {
		view := struct {
			Thead []string
			Tbody []struct {
				Id    string
				Desc  string
				Price string
			}
		}{Thead: []string{"编号", "描述", "价格", "操作"}}
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		gId := r.Form.Get("Id")
		if ts, err := db.FindAllTasks(gId); err != nil {
			flog.LogFile.Println(err)
		} else {
			for _, v := range ts {
				view.Tbody = append(view.Tbody, struct{ Id, Desc, Price string }{v.Id.Hex(), v.Desc, strconv.Itoa(v.Price)})
			}
		}
		if err := t.ExecuteTemplate(w, "task-table", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func memberHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.ParseFiles(filepath.Join(htmlPath, "guild", "html", "member-table.tmpl")); err != nil {
		flog.LogFile.Println(err)
	} else {
		type oneview struct {
			Id      string
			Name    string
			Mobile  string
			Ability string
		}
		view := struct {
			Thead []string
			Tbody []oneview
		}{Thead: []string{"姓名", "手机号", "能力", "操作"}}
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		gId := r.Form.Get("Id")
		if ms, err := db.FindAllMembers(gId); err != nil {
			flog.LogFile.Println(err)
		} else {
			for _, v := range ms {
				view.Tbody = append(view.Tbody, oneview{v.Id.Hex(), v.Name, v.Mobile, v.Ability})
			}
		}
		if err := t.ExecuteTemplate(w, "member-table", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
