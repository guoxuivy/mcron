package mcron

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}
type loginController struct {
}

func (this *loginController) IndexAction(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(TMP_DIR + "/html/login/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

//登陆提交
func (this *loginController) LoginAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 0, "参数错误", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_password := r.FormValue("admin_password")

	if admin_name != "admin" || admin_password != "123456" {
		OutputJson(w, 0, "账户或密码错误", nil)
		return
	}

	// 存入cookie,使用cookie存储
	cookie := http.Cookie{Name: "admin_name", Value: admin_name, Path: "/"}
	http.SetCookie(w, &cookie)
	OutputJson(w, 1, "操作成功", nil)
	return
}

//ajax 返回
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}
