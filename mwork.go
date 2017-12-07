package main

import (
	"fmt"
	"github.com/guoxuivy/mcron"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic", err)
		}
	}()
	mcron.StartServer() //服务端口默认开启一个 任务处理客户端
	// mcron.StartClient()
}
