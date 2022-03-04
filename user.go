package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	// User Channel
	C chan string
	// user socket
	coon net.Conn

	server *Server
}

// create an User API
func NewUser(conn net.Conn, server *Server) *User {
	UserAddr := conn.RemoteAddr().String()

	user := &User{
		Name: UserAddr,
		Addr: UserAddr,
		C:    make(chan string),
		coon: conn,
	}

	// begin to listen user channel message goroutine.
	go user.ListenMessage()

	return user
}

// Send message to user
func (t *User) sendMsg(msg string) {
	t.coon.Write([]byte(msg))
}

// When user online.
func (t *User) Online() {
	// user online, onlineMap add user
	t.server.mapLock.Lock()
	t.server.OnlineMap[t.Name] = t
	t.server.mapLock.Unlock()

	// boardcast user online
	t.server.Broadcast(t, "online")
}

// When user offline
func (t *User) Offline() {
	// user offline, onlineMap remove user
	t.server.mapLock.Lock()
	delete(t.server.OnlineMap, t.Name)
	t.server.mapLock.Unlock()

	// boardcast user offline
	t.server.Broadcast(t, "offline")
}

// user handles message
func (t *User) DoMessage(msg string) {

	t.server.Broadcast(t, msg)
}

// listen user channel function. when message coming, send to client
func (t *User) ListenMessage() {
	for {
		// Get message from channel
		msg := <-t.C
		// Write message
		t.coon.Write([]byte(msg + "\n"))
	}
}
