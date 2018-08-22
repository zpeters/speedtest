package app
import (
	"net"
	"fmt"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/cmds"
	"github.com/zpeters/speedtest/internal/pkg/server"
)

func GetAllServers() (servers []server.Server) {
	return server.GetAllServers()
}

func GetBestServer() (s server.Server) {
	return server.GetBestServer()
}

func Connect(server string) (conn net.Conn) {
	return cmds.Connect(server)
}

func Version(conn net.Conn) (version string) {
	return cmds.Version(conn)
}

func DownloadTest(conn net.Conn, numbytes []int, numtests int) (results float64) {
	var acc float64

	fmt.Printf("Download test: ")
	for i := range numbytes {
		for j := 0; j < numtests; j++ {
			fmt.Printf(".")
			res := cmds.Download(conn, numbytes[i])
			acc = acc + res
		}
	}

	results = acc / float64(numtests)
	fmt.Printf("\n")
	return results
}

func UploadTest(conn net.Conn, numbytes []int, numtests int) (results float64) {
	var acc float64

	fmt.Printf("Upload test: ")
	for i := range numbytes {
		for j := 0; j < numtests; j++ {
			fmt.Printf(".")
			res := cmds.Upload(conn, numbytes[i])
			acc = acc + res
		}
	}

	results = acc / float64(numtests)
	fmt.Printf("\n")
	return results
}


func PingTest(conn net.Conn, numtests int) (results int64) {
	var acc int64

	fmt.Printf("Ping test: ")
	for i := 0; i < numtests; i++ {
		fmt.Printf(".")
		res := cmds.Ping(conn)
		acc = acc + res
	}

	results = acc / int64(numtests)
	fmt.Printf("\n")
	return results
}

