package comms
import (
	"fmt"
	"log"
	"bufio"
	"net"
	"strings"
)

// Public Functions - API
func Connect(server string) (conn net.Conn){
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func Command(conn net.Conn, command string) (resp string) {
	Send(conn, command)
	resp = strings.TrimSpace(string(Recv(conn)))
	return resp
}

func Send(conn net.Conn, msg string) {
	nm := fmt.Sprintf("%s\n", msg)
	//log.Printf("[COMM Tx (%d bytes)] %#v", len(nm), nm)
	fmt.Fprint(conn, nm)
}

func Recv(conn net.Conn) (status []byte) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("[COMM Rx (%d bytes)] %#v", len(data), string(data))

	return data
}

