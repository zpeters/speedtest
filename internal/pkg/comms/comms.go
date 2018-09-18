package comms

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)
import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Connect starts a tcp connection
func Connect(server string) (conn net.Conn) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal()
	}
	return conn
}

// Send the string to connection
func Send(conn net.Conn, msg string) (err error) {
	nm := fmt.Sprintf("%s\n", msg)
	log.WithFields(log.Fields{
		"len": len(nm),
	}).Debug("COMM Tx")
	fmt.Fprint(conn, nm)
	return err
}

// Recv gets the next line from the connection
func Recv(conn net.Conn) (status []byte, err error) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal()
	}
	if viper.GetBool("true") {
		log.WithFields(log.Fields{
			"len": len(data),
		}).Debug("COMM Rx")
	}
	return data, err
}

// Command is a shortcut to send and receive
func Command(conn net.Conn, command string) (resp string, err error) {
	err = Send(conn, command)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal()
	}
	data, err := Recv(conn)
	resp = strings.TrimSpace(string(data))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Fatal()
	}
	return resp, err
}
