package web

import (
	"encoding/json"
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
	model *jobModel
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t := AdminTpl("index")
	list := getModel().getList()
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

//重置任务
func (this *adminController) ReloadAction(w http.ResponseWriter, r *http.Request, user string) {
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 400, "参数错误", nil)
		return
	}
	id := r.FormValue("id")
	jobChan["reload"] <- id
	OutputJson(w, 200, "已重置此任务", nil)
	return
}

//添加任务 (添加逻辑得重做 直接写库 通知chan)
func (this *adminController) AddAction(w http.ResponseWriter, r *http.Request, user string) {
	_id := r.FormValue("scheduleExpr")
	if _id == "" {
		t := AdminTpl("add")
		t.Execute(w, nil)
	} else {
		scheduleExpr := r.FormValue("scheduleExpr")
		desc := r.FormValue("desc")
		shell := r.FormValue("shell")
		ip := r.FormValue("ip")

		job := Job{0, scheduleExpr, desc, shell, ip}
		id, err := getModel().add(job)
		if err != nil {
			OutputJson(w, 400, "添加失败！", err)
		} else {
			jobChan["reload"] <- strconv.Itoa(id) //加载新任务
			OutputJson(w, 200, "添加成功！", nil)
		}
	}
}

//获取一条数据
func (this *adminController) OneAction(w http.ResponseWriter, r *http.Request, user string) {
	id_str := r.FormValue("id")
	id, _ := strconv.Atoi(id_str)
	model := getModel()
	row := model.getOne(id)
	if row.Id > 0 {
		OutputJson(w, 200, "ok", row)
	} else {
		OutputJson(w, 404, "无此数据！", nil)
	}
}

//编辑一条数据
func (this *adminController) EditAction(w http.ResponseWriter, r *http.Request, user string) {
	scheduleExpr := r.FormValue("scheduleExpr")
	desc := r.FormValue("desc")
	shell := r.FormValue("shell")
	ip := r.FormValue("ip")
	id_str := r.FormValue("id")
	id, _ := strconv.Atoi(id_str)

	job := Job{id, scheduleExpr, desc, shell, ip}
	err := getModel().edit(job)
	if err != nil {
		OutputJson(w, 400, "更新失败！", err)
	} else {
		OutputJson(w, 200, "更新成功！", nil)
	}

}
