package mcron

import (
	//"io/ioutil"
	"encoding/json"
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
	clientIp := "127.0.0.1" //默认客户端ip
	conn, err := net.DialTimeout("tcp", clientIp+":"+strconv.Itoa(C_PORT), time.Second*2)
	if err != nil {
		log.Println("连接客户端端失败:", err.Error())
		return
	}
	defer conn.Close()
	shell := job.Shell
	conn.Write([]byte(shell))
	log.Println("server:向客户端发送任务成功：任务ID", job.Id, shell)
}

/**
 * 处理job执行返回结果
 * @param  res string
 * @return bool
 */
func (this *scheduleWorker) backJob(json map[string]interface{}) {
	//time.Sleep(time.Second * 1)
	res := this.getData(json)
	log.Println("server:收到任务反馈数据:", res)
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
		//log.Println("接受客户端连接:", conn.RemoteAddr().String())
		go this.handleConn(conn)
	}
}

func (this *scheduleWorker) handleConn(conn net.Conn) {
	for {
		var buf = make([]byte, 65536)
		n, err := conn.Read(buf)
		if err != nil {
			//log.Println("read error:", err) //连接断开
			break
		}
		if n > 0 {
			json := this.jsonDecode(buf[:n])
			action := this.getAction(json)
			if action == "register" {
				log.Println("server:收到注册信息:", json["Param"])
				conn.Write([]byte("ok"))
			}
			if action == "jobbcak" {
				this.backJob(json)
			}
			//this.backJob(action)

		}
	}
}

//统一传输协议
//b := []byte(`{"Action":"register","Param":["Gomez","Morticia"]}`)
func (this *scheduleWorker) jsonDecode(b []byte) map[string]interface{} {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Println("非json数据：", err)
	}
	m := f.(map[string]interface{})
	return m
	//	log.Println(m["Action"], m["Param"])
	//	for k, v := range m {
	//		switch vv := v.(type) {
	//		case string:
	//			log.Println(k, "is string", vv)
	//		case int:
	//			log.Println(k, "is int", vv)
	//		case []interface{}:
	//			log.Println(k, "is an array:")
	//			for i, u := range vv {
	//				log.Println(i, u)
	//			}
	//		default:
	//			log.Println(k, "is of a type I don't know how to handle")
	//		}
	//	}
}
func (this *scheduleWorker) getAction(m map[string]interface{}) string {
	return this._string("Action", m)
}
func (this *scheduleWorker) getData(m map[string]interface{}) string {
	return this._string("Data", m)
}
func (this *scheduleWorker) _string(key string, m map[string]interface{}) string {
	if nil == m[key] {
		return ""
	} else {
		return m[key].(string)
	}
}
