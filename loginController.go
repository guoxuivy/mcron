package mcron

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type loginController struct {
}

type List []Job

func (this *loginController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TMP_DIR + "/html/login/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

//添加任务
func (this *loginController) AddAction(w http.ResponseWriter, r *http.Request) {
	_id := r.FormValue("id")
	id, _ := strconv.Atoi(_id)
	Server.GetSchedule().AddJob(id, "0/5 * * * * ?")
}

//获取当前任务
func (this *loginController) ListAction(w http.ResponseWriter, r *http.Request) {
	list := Server.GetSchedule().GetJobs()
	//log.Println(list)
	t, err := template.ParseFiles(TMP_DIR + "/html/login/list.html")
	if err != nil {
		log.Println(err)
	}
	data := map[string]List{"List": list}
	t.Execute(w, data)
}
