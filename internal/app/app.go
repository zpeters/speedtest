package app
import (
	"net"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/cmds"
)

func Connect(server string) (conn net.Conn) {
	return cmds.Connect(server)
}

func Version(conn net.Conn) (version string) {
	return cmds.Version(conn)
}

func Quit(conn net.Conn) {
	cmds.Quit(conn)
}

func PingTest(conn net.Conn, numtests int) (results int64) {
	var acc int64

	for i := 0; i < numtests; i++ {
		res := cmds.Ping(conn)
		acc = acc + res
	}

	results = acc / int64(numtests)
	return results
}
