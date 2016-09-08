package mcron

import (
	"net/http"
	"strconv"
)

type User struct {
	UserName string
}

type Page struct {
	UserName string
	List     map[int]Job
}

type adminController struct {
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t := AdminTpl("index")
	list := Server.GetSchedule().GetJobs()
	t.Execute(w, &Page{user, list})
}

//添加任务
func (this *adminController) AddAction(w http.ResponseWriter, r *http.Request, user string) {
	_id := r.FormValue("id")
	if _id == "" {
		t := AdminTpl("add")
		t.Execute(w, nil)
	} else {
		id, _ := strconv.Atoi(_id)
		scheduleExpr := r.FormValue("scheduleExpr")
		desc := r.FormValue("desc")
		msg := Server.GetSchedule().AddJob(id, scheduleExpr, desc)
		OutputJson(w, 0, msg, nil)
	}

}
