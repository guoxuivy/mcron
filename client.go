package mcron

import (
	"io/ioutil"
	"log"
	"net"
	"time"
)

//客户端配置
const (
	//服务器地址
	serverHost = "192.168.37.146:3333"
	//客户端使用端口
	port = 4444
)

//客户端程序
type ClientClass struct {
}

//客户端开启流程
func (this *ClientClass) run() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		log.Println("监听端口失败:", err.Error())
		return
	}
	log.Println("客户端连接已初始化，等待调度指令...")
	this.Listen(listen)
}

//监听来自调度服务器的指令
func (this *ClientClass) Listen(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client:接受客户端连接异常:", err.Error())
			continue
		}
		log.Println("client:收到调度服务器指令:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			//data := make([]byte, 1024)
			result, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Println("读取指令数据错误:", err.Error())
				return
			}
			log.Println("client:收到服务器指令数据:", string(result))
			this.Worker(string(result))
			//conn.Write([]byte(msg))
		}()
	}
}

//向服务器发送数据
func (this *ClientClass) _sendMsg(desc string) {
	//读取客户端配置id
	//conn, err := net.Dial("tcp", "127.0.0.1:3333")
	conn, err := net.Dial("tcp", serverHost)
	if err != nil {
		log.Println("连接服务端端失败:", err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte(desc))
	log.Println("client:处理任务完成：" + desc)
}

//处理指令 返回处理结果
func (this *ClientClass) Worker(shell string) {
	time.Sleep(time.Second * 1)
	this._sendMsg("done")
}

var Client *ClientClass

//创建服务器
func StartClient() {
	Client = &ClientClass{}
	Client.run()
}
