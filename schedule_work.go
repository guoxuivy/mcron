package mcron

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	DSN string = "root:root@tcp(127.0.0.1:3306)/mwork?charset=utf8"
)

/**
 * 数据库连接
 */
func getDb() (*sql.DB, error) {
	db, _ := sql.Open("mysql", DSN)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)
	err := db.Ping()
	if err != nil {
		err = errors.New("数据库连接错误," + fmt.Sprint(DSN))
		return nil, err
	} else {
		return db, nil
	}
}

//任务调度执行者
type scheduleWorker struct {
	db *sql.DB
}

func NewscheduleWorker() (*scheduleWorker, error) {
	DB, err := getDb()
	if nil != err {
		fmt.Println("数据库连接错误！")
		return nil, err
	}
	worker := &scheduleWorker{
		db: DB,
	}
	return worker, nil
}

//每秒一个协程 开始单帧内的工作
func (this *scheduleWorker) Work() {
	begin_time := time.Now().UnixNano()
	fmt.Println(begin_time)
	//读取数据库、redis，将任务分发到client执行，收集结果

}
