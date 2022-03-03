package user

import "net"

type User struct {
	Name string
	Addr string
	// User Channel
	C chan string
	// user socket
	coon net.Conn
}

// create an User API
func NewUser(conn net.Conn) *User {
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

// When user online.
func (t *User) Online() {

}

// When user offline
func (t *User) Offline() {

}

// user handles message
func (t *User) DoMessage() {

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
