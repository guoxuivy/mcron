package mcron

import (
	"errors"
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
	//服务器地址"192.168.38.70" "192.168.48.88"
	SERVER_IP = "127.0.0.1"
)

//客户端程序
type ClientClass struct {
}

//客户端开启流程
func (this *ClientClass) run() {
	err := this.register()
	if nil != err {
		return
	}
	this.Listen()
}

//客户端注册流程
func (this *ClientClass) register() error {
	time.Sleep(time.Millisecond * 100) //0.1s 等待服务端先启动
	conn, err := net.DialTimeout("tcp", SERVER_IP+":"+strconv.Itoa(S_PORT), time.Second*3)
	if err != nil {
		log.Println("注册失败:", err.Error())
		return err
	}
	defer conn.Close()
	conn.Write([]byte(`{"Action":"register","Param":["Param1","Param2"]}`))
	var buf = make([]byte, 65536)
	n, _ := conn.Read(buf)
	json := jsonDecode(buf[:n])
	action := jsonStr("Action", json)
	if action == "register_back" {
		res := jsonStr("Data", json)
		log.Println("client:注册返回:", res)
		if res != "ok" {
			err = errors.New("注册失败")
			return err
		}
	}
	return nil
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
		defer conn.Close()
		go func() {
			result, err := ioutil.ReadAll(conn)
			if err != nil {
				log.Println("读取指令数据错误:", err.Error())
				return
			}

			json := jsonDecode(result)
			act := jsonStr("Action", json)
			if act == "job_run" {
				arr := jsonArr("Data", json)
				id := arr[0].(string)
				shell := arr[1].(string)
				log.Println("client:收到服务器指令数据:", shell)
				this.Worker(id, shell)
			}

		}()
	}
}

//向服务器发送数据
func (this *ClientClass) _sendMsg(id string, desc string) {
	conn, err := net.Dial("tcp", SERVER_IP+":"+strconv.Itoa(S_PORT))
	if err != nil {
		log.Println("连接服务端端失败:", err.Error())
		return
	}
	defer conn.Close()
	conn.Write([]byte(`{"Action":"job_bcak","Data":["` + id + `","` + desc + `"]}`))
}

//处理指令 返回处理结果 服务器记录执行结果日志
func (this *ClientClass) Worker(id string, shell string) {
	log, err := this._execCommand(shell)
	if err != nil {
		go this.WriteLog("shell_run", "["+time.Now().Format("2006-01-02 15:04:05")+"] run ["+shell+"] [error] out: "+err.Error())
		this._sendMsg(id, "error")
	} else {
		go this.WriteLog("shell_run", "["+time.Now().Format("2006-01-02 15:04:05")+"] run ["+shell+"] out: "+log)
		this._sendMsg(id, "done")
	}
}

/**
 * 执行系统命令封装
 * 多个参数以空格分割
 * execCommand("ping baidu.com -n 3")
 */
func (this *ClientClass) _execCommand(shell string) (string, error) {
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
		return "", err
	}
	return string(out), nil
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
