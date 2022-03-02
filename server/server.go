package server

import (
	"fmt"
	"imSystem/user"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	// online user map & map lock
	OnlineMap map[string]*user.User
	mapLock   sync.RWMutex

	// broadcast channel
	Message chan string
}

// Init server API
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user.User),
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
func (t *Server) Broadcast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	fmt.Println(sendMsg)
	t.Message <- sendMsg
}

// bussiness
func (t *Server) Handler(conn net.Conn) {
	// do something
	// fmt.Println("target handler!")

	user := user.NewUser(conn)

	// user online, onlineMap add user
	t.mapLock.Lock()
	t.OnlineMap[user.Name] = user
	t.mapLock.Unlock()

	// boardcast user online
	t.Broadcast(user, "Online")

	// don't let handler die
	select {}
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
