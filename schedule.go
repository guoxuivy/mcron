package mcron

import (
	"io/ioutil"
	"log"
	"net"
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

type CurrJob map[int]Job

//任务管理
type ScheduleManager struct {
	currentJobs CurrJob //当前正在执行的任务id列表
	//周期任务驱动型任务
	cronJob *cron.Cron
	//事件任务驱动型日任务
}

func NewScheduleManager() *ScheduleManager {
	instance := &ScheduleManager{}
	instance.cronJob = cron.New()
	instance.currentJobs = make(map[int]Job)
	return instance
}

func (this *ScheduleManager) Start() {
	this.cronJob.Start()
	//开启客户端监听
	go this._clientListen()
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
}

//任务分发执行
func (this *ScheduleManager) _scheduleActive(id int) {
	log.Println("任务开始执行-任务ID", id)
	job := this.currentJobs[id]
	//根据任务配置分发到相应客户端执行
	//使用tcp通信
	this._sendMsg(job.Desc)
}

//任务分发执行
func (this *ScheduleManager) _sendMsg(desc string) {
	//读取客户端配置id
	conn, err := net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		log.Println("连接客户端端失败:", err.Error())
		return
	}
	defer conn.Close()
	daytime := time.Now().String() + desc
	conn.Write([]byte(daytime))
	log.Println("向客户端发送数据成功：" + daytime)
}

//接受客户端任务反馈
func (this *ScheduleManager) _clientListen() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), 3333, ""})
	if err != nil {
		log.Println("监听端口失败:", err.Error())
		return
	}
	log.Println("已初始化连接，等待客户端反馈...")
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("接受客户端连接异常:", err.Error())
			continue
		}
		//log.Println("收到客户端反馈:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			result, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Println("读取客户端返回数据错误:", err.Error())
				return
			}
			this.Worker(string(result))

		}()
	}
}

//处理指令 返回处理结果
func (this *ScheduleManager) Worker(res string) {
	time.Sleep(time.Second * 1)
	log.Println(" 收到——————任务反馈数据:", res)
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
