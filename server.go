package mcron

import (
	"github.com/guoxuivy/mcron/schedule"
	"github.com/guoxuivy/mcron/webserver"
)

//服务端程序
type Server struct {
	//任务调度中心
	schedule *schedule.ScheduleManager
}

//服务开启流程
func (this *Server) run() {
	this.schedule.Run() //开启任务调度
	webserver.WebRun()  //开启web服务
}

//创建服务器
func StartServer() {
	server := &Server{
		schedule: schedule.NewScheduleManager(),
	}
	server.run()
}
