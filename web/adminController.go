package web

import (
	"encoding/json"
	"net/http"
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
	model := &jobModel{}
	list := model.getList()
	//获取任务运行状态
	jobChan["job_search"] <- "1"
	jobstr := <-jobChan["job_list"]
	ids := []int{}
	if err := json.Unmarshal([]byte(jobstr), &ids); err == nil {
		for _, id := range ids {
			job := list[id]
			job.Desc = list[id].Desc + " run"
			list[id] = job
		}
	}
	t.Execute(w, &Page{user, list})
}

//添加任务
func (this *adminController) AddAction(w http.ResponseWriter, r *http.Request, user string) {
	_id := r.FormValue("scheduleExpr")
	if _id == "" {
		t := AdminTpl("add")
		t.Execute(w, nil)
	} else {
		scheduleExpr := r.FormValue("scheduleExpr")
		desc := r.FormValue("desc")
		shell := r.FormValue("shell")

		//msg := Server.GetSchedule().AddJob(id, scheduleExpr, desc)
		job := &Job{0, scheduleExpr, desc, shell}
		if b, err := json.Marshal(job); err == nil {
			str := string(b)
			jobChan["add"] <- str
		}
		msg := "ok"
		OutputJson(w, 0, msg, nil)
	}

}
