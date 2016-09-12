package mcron

import (
	"log"
	"strconv"
	"time"

	"github.com/guoxuivy/cron"
)

//任务描述
type Job struct {
	Id           int
	ScheduleExpr string
	Desc         string
}

type CurrJob map[int]Job

//任务管理
type ScheduleManager struct {
	currentJobs CurrJob //当前正在执行的任务id列表
	//周期任务驱动型任务
	cronJob *cron.Cron
	//事件任务驱动型日任务
	sWorker *scheduleWorker
}

func NewScheduleManager() *ScheduleManager {
	instance := &ScheduleManager{}
	instance.cronJob = cron.New()
	instance.currentJobs = make(map[int]Job)
	instance.sWorker = NewscheduleWorker()
	return instance
}

func (this *ScheduleManager) Start() {
	//开启定时任务服务
	this.cronJob.Start()
}

func (this *ScheduleManager) Stop() {
	this.cronJob.Stop()
}

//添加任务 任务表达式
func (this *ScheduleManager) GetJobs() CurrJob {
	return this.currentJobs
}

//添加任务 任务表达式
func (this *ScheduleManager) AddJob(id int, scheduleExpr string, desc string) (msg string) {
	//检测是否已存在此id
	if _, ok := this.currentJobs[id]; ok {
		//存在
		return "error"
	}
	job := NewScheduleJob(id, this._scheduleActive)
	this.cronJob.AddJob(scheduleExpr, job, strconv.Itoa(id))
	this.currentJobs[id] = Job{id, scheduleExpr, desc}
	//this.currentJobs = append(this.currentJobs, Job{id, scheduleExpr, desc})
	return "ok"
}

func (this *ScheduleManager) RemoveJob(id int) {
	this.cronJob.RemoveJob(strconv.Itoa(id))
	delete(this.currentJobs, id)
}

//任务分发执行
func (this *ScheduleManager) _scheduleActive(id int) {
	log.Println("server:任务开始执行-任务ID ******start****", id)
	job := this.currentJobs[id]
	go this.sWorker.sendJob(job)
}

//节点、配置心跳监听（待实现 或者直接使用zookeeper）
func (this *ScheduleManager) Monitor() {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
		}
	}()

}

func (this *ScheduleManager) Run() {
	this.Monitor() //异步函数
	this.Start()
	this.sWorker.Start()
	log.Println("任务调度开启")
}
