package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int

	// online user map & map lock
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// broadcast channel
	Message chan string
}

// Init server API
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// listener message broadcast channel's goroutine.
// If have message then send to all user
func (t *Server) ListenMessager() {
	for {
		msg := <-t.Message

		// send message to all online user
		t.mapLock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}
}

// Broadcast function
func (t *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	// show the message on the server console.
	fmt.Println(sendMsg)
	// send message to server channel
	t.Message <- sendMsg
}

// bussiness
func (t *Server) Handler(conn net.Conn) {
	user := NewUser(conn, t)

	// user online
	user.Online()

	// listen is live?
	isLive := make(chan bool)

	// Get Client message (like objective-C 'block function')
	go func() {
		// notice: when user send message is bigger than 4096 btyes, it will break.
		buf := make([]byte, 4096)

		for {
			n, err := conn.Read(buf)

			// Get user send 0, is mean user will offline
			if n == 0 {
				user.Offline()
				return
			}

			// if err is not nil and err is not end of file, is mean Connection read error.
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// get user message (remove '\n')
			msg := string(buf[:n-1])

			// user do somethine for message
			user.DoMessage(msg)

			// user is live
			isLive <- true
		}
	}()

	// don't let handler die
	for {
		select {
		case <-isLive:
			// is live, rebot timer
		case <-time.After(time.Hour * 1):
			// has been over time

			// close user
		}
	}
}

// Start API
func (t *Server) Start() {
	// socket Listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.IP, t.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}

	// close listen
	defer listener.Close()

	// begin message listener's goroutine
	go t.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listerner accept err:", err)
			continue
		}

		// do Handler
		go t.Handler(conn)
	}
}
