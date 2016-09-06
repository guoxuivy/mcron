package mcron

import (
	"log"
	"strconv"
	"time"

	"github.com/jakecoffman/cron"
)

//任务描述
type Job struct {
	Id           int
	ScheduleExpr string
	Desc         string
}

//任务管理
type ScheduleManager struct {
	currentJobs []Job //当前正在执行的任务id列表
	//周期任务驱动型任务
	cronJob *cron.Cron
	//事件任务驱动型日任务
}

func NewScheduleManager() *ScheduleManager {
	instance := &ScheduleManager{}
	instance.cronJob = cron.New()
	//instance.currentJobs = make([]Job, 1)
	return instance
}

func (this *ScheduleManager) Start() {
	this.cronJob.Start()
}

func (this *ScheduleManager) Stop() {
	this.cronJob.Stop()
}

//添加任务 任务表达式
func (this *ScheduleManager) GetJobs() []Job {
	return this.currentJobs
}

//添加任务 任务表达式
func (this *ScheduleManager) AddJob(id int, scheduleExpr string) {
	job := NewScheduleJob(id, this._scheduleActive)
	this.cronJob.AddJob(scheduleExpr, job, strconv.Itoa(id))
	this.currentJobs = append(this.currentJobs, Job{id, scheduleExpr, "test"})
}

func (this *ScheduleManager) RemoveJob(id int) {
	this.cronJob.RemoveJob(strconv.Itoa(id))
}

//需要执行的任务体 php
func (this *ScheduleManager) _scheduleActivePhp(id int) {
	log.Println("php任务开始执行-任务ID", id)
}

//需要执行的任务体 shell
func (this *ScheduleManager) _scheduleActive(id int) {
	log.Println("shell任务开始执行-任务ID", id)
}

//监听是否有任务变化
func (this *ScheduleManager) monitor() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
	}
}

func (this *ScheduleManager) Run() {
	go this.monitor()
	this.Start()
	log.Println("任务调度开启")
	//this.AddJob(2, "0/5 * * * * ?")
}
