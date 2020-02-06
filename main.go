package main

import (
	"fmt"
	"http-mock-server/httpServer"
)

func main() {
	if err := httpServer.Run(); err!=nil{
		fmt.Println(err)
	}
}
