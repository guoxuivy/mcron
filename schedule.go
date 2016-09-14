package mcron

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/guoxuivy/cron"
)

type CurrJob map[int]Job

//任务管理
type ScheduleManager struct {
	addJobChan    chan string //添加任务管道采用json字符串
	removeJobChan chan string //删除任务管道采用json字符串（待实现）
	stopJobChan   chan string //暂停任务管道采用json字符串（待实现）
	startJobChan  chan string //开启任务管道采用json字符串（待实现）

	currentJobs CurrJob         //当前正在执行的任务id列表
	cronJob     *cron.Cron      //周期任务驱动型任务
	sWorker     *scheduleWorker //事件任务驱动型日任务
	jobModel    *jobModel       //数据库操作类
}

func NewScheduleManager() *ScheduleManager {
	instance := &ScheduleManager{}
	instance.addJobChan = make(chan string, 10)
	instance.cronJob = cron.New()
	instance.currentJobs = make(map[int]Job)
	instance.sWorker = &scheduleWorker{}
	instance.jobModel = &jobModel{}
	return instance
}

//启动执行
func (this *ScheduleManager) Start() {
	//开启定时任务服务
	this.cronJob.Start()
	//加载数据库任务
	list := this.jobModel.getList()
	for _, job := range list {
		this._addJob(job)
	}
}

func (this *ScheduleManager) Stop() {
	this.cronJob.Stop()
}

//添加任务 任务表达式
func (this *ScheduleManager) GetJobs() CurrJob {
	return this.currentJobs
}

//写库添加任务
func (this *ScheduleManager) AddJob(j Job) {
	//写入数据库
	id, err := this.jobModel.Add(j)
	if err != nil {
		log.Println(err.Error())
		return
	}
	j.Id = id
	this._addJob(j)
}

//添加任务
func (this *ScheduleManager) _addJob(j Job) {
	job := NewScheduleJob(j.Id, this._scheduleActive)
	this.cronJob.AddJob(j.ScheduleExpr, job, strconv.Itoa(j.Id))
	this.currentJobs[j.Id] = j
}

//零时移除一个执行中的任务（不删除数据库）
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

	//任务添加管道监听
	go func() {
		for {
			var job Job
			jobstr := <-this.addJobChan
			if err := json.Unmarshal([]byte(jobstr), &job); err == nil {
				this.AddJob(job)
			}
		}
	}()
}

func (this *ScheduleManager) Run() {
	this.Monitor() //异步函数
	this.Start()
	this.sWorker.Start()
	log.Println("任务调度开启")
}
