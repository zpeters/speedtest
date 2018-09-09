package comms

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestConnect(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")

	is.True(conn != nil) // conn should be something
}

func TestSend(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	err := Send(conn, "PING")

	is.NoErr(err) // should be able to send something
}

func TestRecv(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	err := Send(conn, "HI")
	data, err2 := Recv(conn)

	is.NoErr(err)                                    // should be able to send
	is.NoErr(err2)                                   // should be able to recieve
	is.True(strings.Contains(string(data), "HELLO")) // we should have HELLO back
}

func TestCommand(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	resp, err := Command(conn, "PING 123")

	is.NoErr(err)                                   // command should return a response
	is.True(strings.Contains(string(resp), "PONG")) // PING should get a PONG
}
