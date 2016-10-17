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
	mcron.StartServer()
}
