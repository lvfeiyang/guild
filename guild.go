package main

import (
	"github.com/lvfeiyang/proxy/common/config"
	"github.com/lvfeiyang/proxy/common/db"
	gDb "github.com/lvfeiyang/guild/common/db"
	"github.com/lvfeiyang/proxy/common/flog"
	"github.com/lvfeiyang/proxy/common/session"
	"github.com/lvfeiyang/proxy/common"
	"github.com/lvfeiyang/guild/message"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

var htmlPath string
var pjtCfg config.ProjectConfig

func main() {
	flog.Init()
	config.Init()
	session.Init()
	db.Init()
	message.Init()
	httpAddr := ":80"
	htmlPath = config.ConfigVal.HtmlPath
	if pjtCfg = config.GetProjectConfig("guild"); "" == pjtCfg.Name {
		flog.LogFile.Fatal("no guild project!")
	}

	if !pjtCfg.Proxy {
		jsFiles, cssFiles, fontsFiles := filepath.Join(htmlPath, "sfk", "js"), filepath.Join(htmlPath, "sfk", "css"), filepath.Join(htmlPath, "sfk", "fonts")
		http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(jsFiles))))
		http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(cssFiles))))
		http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir(fontsFiles))))

		http.Handle("/guild/msg/", &message.LocMessage{})
	} else {
		httpAddr = pjtCfg.Http
	}

	gcssFiles := filepath.Join(htmlPath, "guild", "html", "css")
	gjsFiles := filepath.Join(htmlPath, "guild", "html", "js")
	http.Handle("/guild/css/", http.StripPrefix("/guild/css/", http.FileServer(http.Dir(gcssFiles))))
	http.Handle("/guild/js/", http.StripPrefix("/guild/js/", http.FileServer(http.Dir(gjsFiles))))

	go common.ListenTcp(pjtCfg.Tcp, message.MhMap)

	http.HandleFunc("/guild", guildHandler)
	http.HandleFunc("/guild/detail", guildDetailHandler)
	http.HandleFunc("/guild/task", taskHandler)
	http.HandleFunc("/guild/member", memberHandler)

	flog.LogFile.Fatal(http.ListenAndServe(httpAddr, nil))
}
func guildHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "guild.html"),
		filepath.Join(htmlPath, "guild", "html", "sidebar.tmpl"),

		filepath.Join(htmlPath, "guild", "html", "modal", "edit-guild.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "edit-member.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "edit-task.tmpl"),
		filepath.Join(htmlPath, "guild", "html", "modal", "login.tmpl"),
	}

	if t, err := template.ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		var view struct {
			GuildList []struct {
				Id   string
				Name string
			}
		}
		gs, err := gDb.FindAllGuilds()
		if err != nil {
			flog.LogFile.Println(err)
		}
		for _, v := range gs {
			view.GuildList = append(view.GuildList, struct{ Id, Name string }{v.Id.Hex(), v.Name})
		}
		if err := t.Execute(w, view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func haveRole(all, one byte) bool {
	return 0 != all&one
}
func guildDetailHandler(w http.ResponseWriter, r *http.Request) {
	paths := []string{
		filepath.Join(htmlPath, "guild", "html", "main.tmpl"),
	}
	if t, err := template.New("main").Funcs(
		template.FuncMap{"haveRole": haveRole}).ParseFiles(paths...); err != nil {
		flog.LogFile.Println(err)
	} else {
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		id := r.Form.Get("Id")
		g := gDb.Guild{}
		if bson.IsObjectIdHex(id) {
			(&g).GetById(bson.ObjectIdHex(id))
		}
		role, err := gDb.RoleAble(r.Header.Get("SessionId"), id)
		if err != nil {
			flog.LogFile.Println(err)
		}
		view := struct {
			Id        string
			Name      string
			Introduce string
			Role      byte
		}{g.Id.Hex(), g.Name, g.Introduce, role}
		if err := t.ExecuteTemplate(w, "main", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func taskHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.New("task").Funcs(template.FuncMap{"haveRole": haveRole}).
		ParseFiles(filepath.Join(htmlPath, "guild", "html", "task-table.tmpl")); err != nil {
		flog.LogFile.Println(err)
	} else {
		view := struct {
			Thead []string
			Tbody []struct {
				Id    string
				Desc  string
				Price string
			}
			Role byte
		}{Thead: []string{"编号", "描述", "价格", "操作"}}
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		gId := r.Form.Get("Id")
		if ts, err := gDb.FindAllTasks(gId); err != nil {
			flog.LogFile.Println(err)
		} else {
			for _, v := range ts {
				view.Tbody = append(view.Tbody, struct{ Id, Desc, Price string }{v.Id.Hex(), v.Desc, strconv.Itoa(v.Price)})
			}
		}
		role, err := gDb.RoleAble(r.Header.Get("SessionId"), gId)
		if err != nil {
			flog.LogFile.Println(err)
		}
		view.Role = role
		if err := t.ExecuteTemplate(w, "task-table", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
func memberHandler(w http.ResponseWriter, r *http.Request) {
	if t, err := template.New("member").Funcs(template.FuncMap{"haveRole": haveRole}).
		ParseFiles(filepath.Join(htmlPath, "guild", "html", "member-table.tmpl")); err != nil {
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
			Role  byte
		}{Thead: []string{"姓名", "手机号", "能力", "操作"}}
		if err := r.ParseForm(); err != nil {
			flog.LogFile.Println(err)
		}
		gId := r.Form.Get("Id")
		if ms, err := gDb.FindAllMembers(gId); err != nil {
			flog.LogFile.Println(err)
		} else {
			for _, v := range ms {
				view.Tbody = append(view.Tbody, oneview{v.Id.Hex(), v.Name, v.Mobile, v.Ability})
			}
		}
		role, err := gDb.RoleAble(r.Header.Get("SessionId"), gId)
		if err != nil {
			flog.LogFile.Println(err)
		}
		view.Role = role
		if err := t.ExecuteTemplate(w, "member-table", view); err != nil {
			flog.LogFile.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
