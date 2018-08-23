package comms

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// Connect starts a tcp connection
func Connect(server string) (conn net.Conn) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// Send the string to connection
func Send(conn net.Conn, msg string) {
	nm := fmt.Sprintf("%s\n", msg)
	// For DEBUGGING
	//log.Printf("[COMM Tx (%d bytes)] %#v", len(nm), nm)
	fmt.Fprint(conn, nm)
}

// Recv gets the next line from the connection
func Recv(conn net.Conn) (status []byte) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	// For DEBUGGING
	//log.Printf("[COMM Rx (%d bytes)] %#v", len(data), string(data))

	return data
}

// Command is a shortcut to send and receive
func Command(conn net.Conn, command string) (resp string) {
	Send(conn, command)
	resp = strings.TrimSpace(string(Recv(conn)))
	return resp
}
