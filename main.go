package main

import (
	"fmt"
	"imSystem/server"
)

func main() {
	fmt.Println("Start Program.")

	ser := server.NewServer("127.0.0.1", 8888)
	ser.Start()

	fmt.Println("done.")
}
