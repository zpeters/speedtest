package comms

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)
import (
	"github.com/spf13/viper"
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
func Send(conn net.Conn, msg string) (err error) {
	nm := fmt.Sprintf("%s\n", msg)
	if viper.GetBool("Debug") {
		log.Printf("[COMM Tx (%d bytes)] %#v", len(nm), nm)
	}
	fmt.Fprint(conn, nm)
	return err
}

// Recv gets the next line from the connection
func Recv(conn net.Conn) (status []byte, err error) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	if viper.GetBool("true") {
		log.Printf("[COMM Rx (%d bytes)] %#v", len(data), string(data))
	}
	return data, err
}

// Command is a shortcut to send and receive
func Command(conn net.Conn, command string) (resp string, err error) {
	err = Send(conn, command)
	if err != nil {
		log.Fatal(err)
	}
	data, err := Recv(conn)
	resp = strings.TrimSpace(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}
