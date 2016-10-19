package mcron

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/guoxuivy/mcron/web"
)

//任务描述
type Job web.Job

//任务日志
type JobLog struct {
	Id         int
	JobId      int
	Action     string
	Log        string
	CreateTime string
}

type Model struct {
	web.JobModel
}

//添加任务日志
func (this *Model) addLog(log JobLog) (int, error) {
	//写入数据库
	db, err := web.GetDb()
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
