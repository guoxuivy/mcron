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
	currentJobs CurrJob                //当前正在执行的任务id列表
	cronJob     *cron.Cron             //周期任务驱动型任务
	sWorker     *scheduleWorker        //事件任务驱动型日任务
	jobModel    *jobModel              //数据库操作类
	jobChan     map[string]chan string //任务操作管道
	jobLogChan  chan JobLog            //任务日志记录管道
}

func NewScheduleManager() *ScheduleManager {
	chans := jobActiveChan{
		"add":    make(chan string, 10), //json_str
		"remove": make(chan string, 10),
		"stop":   make(chan string, 10),
		"start":  make(chan string, 10),
		"reload": make(chan string, 10),

		"job_search": make(chan string, 1), //web 当前任务查询请求
		"job_list":   make(chan string, 1), //web 当前任务查询结果返回 json 传送
	}
	instance := &ScheduleManager{}
	instance.jobChan = chans
	instance.cronJob = cron.New()
	instance.currentJobs = make(map[int]Job)
	instance.sWorker = &scheduleWorker{}
	instance.jobModel = &jobModel{}
	instance.jobLogChan = make(chan JobLog, 1000)

	return instance
}

//启动执行
func (this *ScheduleManager) Start() {
	//开启定时任务服务
	this.cronJob.Start()
	//加载数据库任务
	list, err := this.jobModel.getList()
	if err != nil {
		panic("服务器启动失败：" + err.Error())
	}
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
//func (this *ScheduleManager) AddJob(j Job) {
//	id, err := this.jobModel.Add(j)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	j.Id = id
//	this._addJob(j)
//}

//从数据库重载任务
func (this *ScheduleManager) ReloadJob(id int) {
	job := this.jobModel.getOne(id)
	this.RemoveJob(id)
	this._addJob(job)
}

//删除任务
func (this *ScheduleManager) DeleteJob(id int) {
	//先停止
	this.RemoveJob(id)
	//再数据库删除

}

//零时移除一个执行中的任务（不删除数据库） stop
func (this *ScheduleManager) RemoveJob(id int) {
	this.cronJob.RemoveJob(strconv.Itoa(id))
	delete(this.currentJobs, id)
}

//添加任务
func (this *ScheduleManager) _addJob(j Job) {
	job := NewScheduleJob(j.Id, this._scheduleActive)
	this.cronJob.AddJob(j.ScheduleExpr, job, strconv.Itoa(j.Id))
	this.currentJobs[j.Id] = j
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
	go this.doLog()

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
			case <-this.jobChan["add"]: //已废弃
				//				var job Job
				//				if err := json.Unmarshal([]byte(jobstr), &job); err == nil {
				//					this.AddJob(job)
				//				}
			case jobid := <-this.jobChan["remove"]: //彻底删除任务
				id, _ := strconv.Atoi(jobid)
				this.DeleteJob(id)
				log.Println("任务删除：", jobid)
			case jobid := <-this.jobChan["stop"]: //暂停任务
				id, _ := strconv.Atoi(jobid)
				this.RemoveJob(id)
				log.Println("任务暂停：", jobid)
				this.WriteLog(id, "job_stop", "任务被暂停")
			case jobid := <-this.jobChan["start"]: //开启暂停中的任务
				log.Println(jobid)
			case jobid := <-this.jobChan["reload"]: //彻底删除任务
				id, _ := strconv.Atoi(jobid)
				this.ReloadJob(id)
				log.Println("任务重载：", jobid)
				this.WriteLog(id, "job_reload", "任务被重载")
			case jobid := <-this.jobChan["job_search"]: //查询运行时任务
				id, _ := strconv.Atoi(jobid)
				if id > 0 {
					if _, ok := this.currentJobs[id]; ok {
						//存在
						this.jobChan["job_list"] <- "1"
					} else {
						this.jobChan["job_list"] <- "0"
					}
				} else {
					ids := []int{}
					for id, _ := range this.currentJobs {
						ids = append(ids, id)
					}
					if b, err := json.Marshal(ids); err == nil {
						this.jobChan["job_list"] <- string(b)
					}
				}
			}
		}
	}()
}

/**
 * 任务日志写入队列
 * @param id 任务id
 * @param act 任务操作
 * @param log 具体日志
 * @return
 */
func (this *ScheduleManager) WriteLog(jobid int, act string, logstr string) {
	var log JobLog
	log.JobId = jobid
	log.Action = act
	log.Log = logstr
	this.jobLogChan <- log
}

/**
 * 任务日志入库处理
 */
func (this *ScheduleManager) doLog() {
	for {
		jobLog := <-this.jobLogChan
		_, err := this.jobModel.AddLog(jobLog)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (this *ScheduleManager) Run() {
	this.Monitor() //异步函数
	this.Start()
	this.sWorker.Start()
	log.Println("任务调度开启")
}
