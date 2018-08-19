package app
import (
	"net"
	"fmt"
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


func DownloadTest(conn net.Conn, numtests int, bytes int) (results string) {
	var acc float64

	for i := 0; i < numtests; i++ {
		res := cmds.Download(conn, bytes)
		acc = acc + res
		fmt.Printf("Download Test %d %f\n", i, res)
	}

	resFloat := acc / float64(numtests)
	results = fmt.Sprintf("%f", resFloat)
	return results
}

func PingTest(conn net.Conn, numtests int) (results int64) {
	var acc int64

	for i := 0; i < numtests; i++ {
		res := cmds.Ping(conn)
		acc = acc + res
		fmt.Printf("Ping Test %d - %d\n", i, res)
	}

	results = acc / int64(numtests)
	return results
}
