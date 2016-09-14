package mcron

import (
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	//服务端port 11510
	//客户端port 11511
	S_PORT int = 11510
	C_PORT int = 11511
)

//任务调度执行者
type scheduleWorker struct {
}

func (this *scheduleWorker) Start() {
	go this._clientListen() //启动job反馈侦听
}

//向客户端发送任务
func (this *scheduleWorker) sendJob(job Job) {
	//根据任务配置分发到相应客户端执行
	//读取客户端配置id
	//conn, err := net.Dial("tcp", "127.0.0.1:4444")
	//clientIp := "192.168.51.125"
	clientIp := "127.0.0.1"
	conn, err := net.DialTimeout("tcp", clientIp+":"+strconv.Itoa(C_PORT), time.Second*2)
	if err != nil {
		log.Println("连接客户端端失败:", err.Error())
		return
	}
	defer conn.Close()
	daytime := time.Now().String() + job.Desc
	conn.Write([]byte(daytime))
	log.Println("server:向客户端发送任务成功：任务ID", job.Id, daytime)
}

//job返回结果处理 需要定义一套标准返回协议
func (this *scheduleWorker) backJob(res string) {
	time.Sleep(time.Second * 1)
	log.Println("server:收到任务反馈数据:", res)
}

//job执行反馈侦听
func (this *scheduleWorker) _clientListen() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), 11510, ""})
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
			this.backJob(string(result))
		}()
	}
}
