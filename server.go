package mcron

//服务端程序
type ServerClass struct {
	//任务调度中心
	schedule *ScheduleManager
}

//服务开启流程
func (this *ServerClass) run() {
	this.schedule.Run() //开启任务调度
	WebRun()            //开启web服务
}

func (this *ServerClass) GetSchedule() *ScheduleManager {
	return this.schedule
}

func (this *ServerClass) ListJob() CurrJob {
	return this.schedule.GetJobs()
}

var Server *ServerClass

//创建服务器
func StartServer() {
	Server = &ServerClass{
		schedule: NewScheduleManager(),
	}
	Server.run()
}
