package main

import (
	"fmt"
	"imSystem/server"
)

// Start the Program.
func ProgramStart() {
	ser := server.NewServer("127.0.0.1", 8888)
	ser.Start()
}

// Program Entrance
func main() {
	fmt.Println("Start Program.")

	ProgramStart()

	fmt.Println("done.")
}
