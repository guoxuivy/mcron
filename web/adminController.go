package web

import (
	"encoding/json"
	"log"
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

//获取joblist
func (this *adminController) getList() map[int]map[string]string {
	db, err := getDb()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	//查询数据库
	query, err := db.Query("SELECT * FROM `job_list` WHERE `status` = ? ", 1)
	if err != nil {
		log.Println("查询数据库失败", err.Error())
		return nil
	}
	defer query.Close()

	//读出查询出的列字段名
	cols, _ := query.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}

	//最后得到的map
	results := make(map[int]map[string]string)
	i := 0
	for query.Next() { //循环，让游标往下推
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			log.Println(err)
			return nil
		}

		row := make(map[string]string) //每行数据

		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}

	return results
}

func (this *adminController) IndexAction(w http.ResponseWriter, r *http.Request, user string) {
	t := AdminTpl("index")
	model := &jobModel{}
	list := model.getList()
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
		//msg := Server.GetSchedule().AddJob(id, scheduleExpr, desc)

		job := &Job{0, scheduleExpr, desc}

		if b, err := json.Marshal(job); err == nil {
			str := string(b)
			jobChan <- str
		}

		msg := "ok"
		OutputJson(w, 0, msg, nil)
	}

}
