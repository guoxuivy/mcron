package mcron

import (
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
}

type adminController struct {
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t, err := template.ParseFiles(TMP_DIR + "/html/admin/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &User{user})
}
