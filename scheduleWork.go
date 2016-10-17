package mcron

import (
	"log"
	"net"
	"strconv"
	"time"
)

const (
	S_PORT int = 11510 //服务端port
	C_PORT int = 11511 //客户端port
)

//任务调度执行者
type scheduleWorker struct {
}

//启动流程
func (this *scheduleWorker) Start() {
	go this._clientListen()
}

/**
 * 向客户端发送任务
 * @param  job Job
 * @return bool
 */
func (this *scheduleWorker) sendJob(job Job) {
	//根据任务配置分发到相应客户端执行
	//clientIp := "127.0.0.1" //默认客户端ip
	clientIp := job.IP
	conn, err := net.DialTimeout("tcp", clientIp+":"+strconv.Itoa(C_PORT), time.Second*2)
	if err != nil {
		Server.Schedule.WriteLog(job.Id, "send_job", "发送任务失败:"+err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte(`{"Action":"job_run","Data":["` + strconv.Itoa(job.Id) + `","` + job.Shell + `"]}`))
}

/**
 * 处理job执行返回结果
 * @param  res string
 * @return bool
 */
func (this *scheduleWorker) backJob(json map[string]interface{}) {
	//time.Sleep(time.Second * 1)
	arr := jsonArr("Data", json)
	idstr := arr[0].(string)
	err := arr[1].(string)
	id, _ := strconv.Atoi(idstr)
	if err == "error" {
		log.Println("任务执行失败 ---", idstr)
		//扩展监控警报处理
	}
	Server.Schedule.WriteLog(id, "job_run_back", err)
}

/**
 * 客户端执行反馈侦听
 * @param
 * @return
 */
func (this *scheduleWorker) _clientListen() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), 11510, ""})
	if err != nil {
		log.Println("监听端口失败:", err.Error())
		return
	}
	defer listen.Close()
	log.Println("服务端已初始化连接，等待客户端反馈...")
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("接受客户端连接异常:", err.Error())
			continue
		}
		defer conn.Close()
		go this.handleConn(conn)
	}
}

func (this *scheduleWorker) handleConn(conn net.Conn) {
	for {
		var buf = make([]byte, 65536)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		if n > 0 {
			json := jsonDecode(buf[:n])
			action := jsonStr("Action", json)
			if action == "register" {
				log.Println("server:收到注册信息:", json["Param"])
				conn.Write([]byte(`{"Action":"register_back","Data":"ok"}`))
			}
			if action == "job_bcak" {
				this.backJob(json)
			}
		}
	}
}
