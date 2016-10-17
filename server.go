package mcron

import (
	"github.com/guoxuivy/mcron/web"
)

//服务端程序
type ServerClass struct {
	//任务调度中心
	Schedule *ScheduleManager
}

//服务开启流程
func (this *ServerClass) Start() {
	go StartClient()                  //开启自身任务处理客户端 如没有可不开启
	this.Schedule.Run()               //开启任务调度
	web.WebRun(this.Schedule.jobChan) //开启web服务 阻塞式
}

func (this *ServerClass) GetSchedule() *ScheduleManager {
	return this.Schedule
}

func (this *ServerClass) ListJob() CurrJob {
	return this.Schedule.GetJobs()
}

var Server *ServerClass

//创建服务器
func StartServer() {
	Server = &ServerClass{
		Schedule: NewScheduleManager(),
	}
	Server.Start()
}
