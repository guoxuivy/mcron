<h2>分布式任务调度系统</h2> 

<b>包括服务端和客户端</b> 

<pre>
服务端：
    1、提供定时任务管理（web端可视化）
    2、任务调度分发
    3、收集执行结果
    4、客户端监听（配置与状态）
客户端：
    1、接受执行指令，执行任务并反馈结果
    2、心跳状态报告
</pre>

模板文件路径需要自行配置：TMP_DIR = "template"

开启实例：
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic", err)
		}
	}()
	mcron.StartServer() //服务端口默认开启一个 任务处理客户端
	// mcron.StartClient()
}