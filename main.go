package main

import (
	"imSystem/server"
)

func main() {
	ser := server.NewServer("127.0.0.1", 8888)
	ser.Start()
}
