package webserver

import (
	"log"
	"net/http"
)

//web服务器开始工作（可以用已有的php替换自带的web服务器）
func WebRun() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/", NotFoundHandler)
	http.ListenAndServe(":8888", nil)
	log.Println("web服务器开启")
}
