package web

import (
	"log"
	"net/http"
)

const (
	TMP_DIR = "github.com/guoxuivy/mcron/template"
)

//任务在线处理管道
var jobChan map[string]chan string

//web服务器开始工作（可以用已有的php替换自带的web服务器）
func WebRun(chans map[string]chan string) {
	jobChan = chans
	http.Handle("/css/", http.FileServer(http.Dir(TMP_DIR)))
	http.Handle("/js/", http.FileServer(http.Dir(TMP_DIR)))

	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/", NotFoundHandler)
	http.ListenAndServe(":8888", nil)
	log.Println("web服务器开启")
}
