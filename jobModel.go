package mcron

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
	IP           string
}

//任务日志
type JobLog struct {
	Id         int
	JobId      int
	Action     string
	Log        string
	CreateTime string
}

type jobModel struct{}

//添加任务 任务表达式
func (this *jobModel) Add(j Job) (int, error) {
	//写入数据库
	db, err := getDb()
	if err != nil {
		return 0, err
	}
	stmt, _ := db.Prepare("INSERT INTO `job_list` (`schedule_expr`, `desc`, `shell`, `ip`, `status`) VALUES (?,?,?,?,1)")
	defer stmt.Close()
	res, err := stmt.Exec(j.ScheduleExpr, j.Desc, j.Shell, j.IP)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

//添加任务日志
func (this *jobModel) AddLog(log JobLog) (int, error) {
	//写入数据库
	db, err := getDb()
	if err != nil {
		return 0, err
	}
	stmt, err := db.Prepare("INSERT INTO `job_log` (`job_id`, `action`, `log`) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(log.JobId, log.Action, log.Log)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), err
}

func (this *jobModel) getList() (map[int]Job, error) {
	jobs := make(map[int]Job)
	res, err := this.findAll()
	if res != nil {
		for _, v := range res {
			id, _ := strconv.Atoi(v["id"])
			jobs[id] = Job{id, v["schedule_expr"], v["desc"], v["shell"], v["ip"]}
		}
	}
	return jobs, err
}

func (this *jobModel) getOne(id int) Job {
	db, err := getDb()
	var job Job
	one, err := db.Query("SELECT `id`, `schedule_expr`, `desc`, `shell`, `ip` FROM `job_list` WHERE `id` = ? ", id)
	if err != nil {
		log.Println(err)
	}
	defer one.Close()
	for one.Next() {
		err := one.Scan(&job.Id, &job.ScheduleExpr, &job.Desc, &job.Shell, &job.IP)
		if err != nil {
			log.Println(err)
		}
	}
	return job
}

//通用列表查询
func (this *jobModel) findAll() (map[int]map[string]string, error) {
	db, err := getDb()
	if err != nil {
		return nil, err
	}
	//查询数据库
	query, err := db.Query("SELECT * FROM `job_list` WHERE `status` = ? ", 1)
	if err != nil {
		return nil, err
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
	for query.Next() {
		//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}

	return results, nil
}
