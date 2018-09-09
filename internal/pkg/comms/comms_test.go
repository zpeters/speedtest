package comms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	assert.NotNil(t, conn)
}

func TestSend(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	err := Send(conn, "PING")
	assert.NoError(t, err)
}

func TestRecv(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	err := Send(conn, "HI")
	assert.NoError(t, err)
	data, err := Recv(conn)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "HELLO")
}

func TestCommand(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	resp, err := Command(conn, "PING 123")
	assert.NoError(t, err)
	assert.Contains(t, string(resp), "PONG")
}
