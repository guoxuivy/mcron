package mcron

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	UserName string
}

type Page struct {
	UserName string
	List     []Job
}

type adminController struct {
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles(TMP_DIR + "/html/admin/index.html")
	if err != nil {
		log.Println(err)
	}

	list := Server.GetSchedule().GetJobs()

	t.Execute(w, &Page{user, list})
}

//添加任务
func (this *adminController) AddAction(w http.ResponseWriter, r *http.Request, user string) {
	_id := r.FormValue("id")
	id, _ := strconv.Atoi(_id)
	Server.GetSchedule().AddJob(id, "0/5 * * * * ?")
}
