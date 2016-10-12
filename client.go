package mcron

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//客户端配置
const (
	//服务器地址 只响应 json格式 必须包含action
	SERVER_IP = "192.168.38.70"
)

//客户端程序
type ClientClass struct {
}

//客户端开启流程
func (this *ClientClass) run() {
	this.register()
	this.Listen()
}

//客户端注册流程
func (this *ClientClass) register() {
	conn, err := net.Dial("tcp", SERVER_IP+":"+strconv.Itoa(S_PORT))
	if err != nil {
		log.Println("连接服务端端失败:", err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte(`{"Action":"register","Param":["Param1","Param2"]}`))

	var buf = make([]byte, 65536)
	n, _ := conn.Read(buf)
	res := string(buf[:n])
	if res == "ok" {
		log.Println("注册成功！")
		return
	}

}

//监听来自调度服务器的指令
func (this *ClientClass) Listen() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), C_PORT, ""})
	if err != nil {
		log.Println("监听端口失败:", err.Error())
		return
	}
	log.Println("客户端连接已初始化，等待调度指令...")
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("client:接受客户端连接异常:", err.Error())
			continue
		}
		//log.Println("client:收到调度服务器指令:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			result, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Println("读取指令数据错误:", err.Error())
				return
			}
			log.Println("client:收到服务器指令数据:", string(result))
			this.Worker(string(result))
		}()
	}
}

//向服务器发送数据
func (this *ClientClass) _sendMsg(desc string) {
	conn, err := net.Dial("tcp", SERVER_IP+":"+strconv.Itoa(S_PORT))
	if err != nil {
		log.Println("连接服务端端失败:", err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte(`{"Action":"jobbcak","Data":"` + desc + `"}`))
	//conn.Write([]byte(desc))
	//log.Println("client:处理任务完成：" + desc)
}

//处理指令 返回处理结果
func (this *ClientClass) Worker(shell string) {
	//time.Sleep(time.Second * 1)
	log := this._execCommand(shell)
	//执行结果日志记录在本地，向调度中心返回执行结果即可
	go this.WriteLog("shell_run", "["+time.Now().Format("2006-01-02 15:04:05")+"] run ["+shell+"] out: "+log)
	this._sendMsg("done")
}

/**
 * 执行系统命令封装
 * 多个参数以空格分割
 * execCommand("ping baidu.com -n 3")
 */
func (this *ClientClass) _execCommand(shell string) string {
	params := strings.Split(shell, " ")
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		params = append([]string{"/C"}, params...)
		cmd = exec.Command("cmd", params...)
	} else {
		command := params[0]
		params = params[1:]
		cmd = exec.Command(command, params...)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

/**
 * 升级日志写入  文件追加
 * @param  {[type]} log string        [description]
 * @return {[type]}     [description]
 */
func (this *ClientClass) WriteLog(tag string, data string) {
	str_time := time.Now().Format("2006_01_02")
	filename := tag + "_" + str_time + ".log"
	fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
	}
	defer fl.Close()
	fl.WriteString(data)
	fl.WriteString("\r\n")
}

var Client *ClientClass

//创建客户端
func StartClient() {
	Client = &ClientClass{}
	Client.run()
}
