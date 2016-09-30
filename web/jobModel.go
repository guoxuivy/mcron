package web

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MYSQL_DSN string = "root:root@tcp(127.0.0.1:3306)/mwork?charset=utf8"
)

var DB *sql.DB

/**
 * 数据库连接
 */
func getDb() (*sql.DB, error) {
	if DB == nil {
		db, _ := sql.Open("mysql", MYSQL_DSN)
		db.SetMaxOpenConns(200)
		db.SetMaxIdleConns(100)
		err := db.Ping()
		if err != nil {
			err = errors.New("数据库连接错误," + MYSQL_DSN)
			return nil, err
		} else {
			DB = db
		}
	}
	return DB, nil
}

//数据库job 映射

//任务描述
type Job struct {
	Id           int
	ScheduleExpr string
	Desc         string
	Shell        string
}

type jobModel struct{}

func (this *jobModel) getList() map[int]Job {
	jobs := make(map[int]Job)
	res := this.findAll()
	if res == nil {
		log.Println("empty")
	} else {
		for _, v := range res {
			id, _ := strconv.Atoi(v["id"])
			jobs[id] = Job{id, v["schedule_expr"], v["desc"], v["shell"]}
		}
	}
	return jobs
}

func (this *jobModel) getShell(id int) string {
	db, err := getDb()
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	var shell string
	err = db.QueryRow("SELECT shell FROM `job_list` WHERE `id`=?", id).Scan(&shell)
	return shell
}

//通用列表查询
func (this *jobModel) findAll() map[int]map[string]string {
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
