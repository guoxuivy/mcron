package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

//ajax标准输出
type Result struct {
	Code int //  200/400
	Msg  string
	Data interface{}
}

//ajax 返回
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	w.Header().Set("content-type", "application/json")
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

//admin获取模版文件
func AdminTpl(name string) *template.Template {
	t, err := template.ParseFiles(TMP_DIR+"/html/admin/"+name+".html", TMP_DIR+"/html/admin/layout.html")
	if err != nil {
		log.Println(err)
	}
	return t
}
