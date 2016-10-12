package mcron

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/guoxuivy/cron"
)

type CurrJob map[int]Job

//任务在线处理管道
type jobActiveChan map[string]chan string

//任务管理
type ScheduleManager struct {
	currentJobs CurrJob         //当前正在执行的任务id列表
	cronJob     *cron.Cron      //周期任务驱动型任务
	sWorker     *scheduleWorker //事件任务驱动型日任务
	jobModel    *jobModel       //数据库操作类
	jobChan     map[string]chan string
}

func NewScheduleManager() *ScheduleManager {
	chans := jobActiveChan{
		"add":    make(chan string, 10),
		"remove": make(chan string, 10),
		"stop":   make(chan string, 10),
		"start":  make(chan string, 10),

		"job_search": make(chan string, 1), //web 当前任务查询请求
		"job_list":   make(chan string, 1), //web 当前任务查询结果返回 json 传送
	}
	instance := &ScheduleManager{}
	instance.jobChan = chans
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

/**
 * 服务端监听 包括（客户端消息、客户端配置、客户端心跳等）或者直接使用zookeeper
 * @param
 * @return
 */
func (this *ScheduleManager) Monitor() {
	go func() {
		//心跳（每秒）
		ticker := time.NewTicker(time.Second)
		for {
			<-ticker.C
		}
	}()

	//web前台任务操作管道监听
	go func() {
		for {
			select {
			case jobstr := <-this.jobChan["add"]:
				var job Job
				if err := json.Unmarshal([]byte(jobstr), &job); err == nil {
					this.AddJob(job)
				}
			case jobid := <-this.jobChan["remove"]:
				log.Println(jobid)
			case jobid := <-this.jobChan["stop"]:
				log.Println(jobid)
			case jobid := <-this.jobChan["start"]:
				log.Println(jobid)
			case <-this.jobChan["job_search"]:
				//非阻塞式
				go func() {
					ids := []int{}
					for id, _ := range this.currentJobs {
						ids = append(ids, id)
					}
					if b, err := json.Marshal(ids); err == nil {
						this.jobChan["job_list"] <- string(b)
					}
				}()
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
