package web

import (
	"html/template"
	"log"
	"net/http"
)

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
	err := r.ParseForm()
	if err != nil {
		OutputJson(w, 400, "参数错误", nil)
		return
	}

	admin_name := r.FormValue("admin_name")
	admin_password := r.FormValue("admin_password")

	if admin_name != "admin" || admin_password != "123456" {
		OutputJson(w, 400, "账户或密码错误", nil)
		return
	}

	// 存入cookie,使用cookie存储
	cookie := http.Cookie{Name: "admin_name", Value: admin_name, Path: "/"}
	http.SetCookie(w, &cookie)
	OutputJson(w, 200, "操作成功", nil)
	return
}
