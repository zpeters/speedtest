package app

import (
	"fmt"
	"net"
	"time"
)
import (
	"github.com/spf13/viper"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/cmds"
	"github.com/zpeters/speedtest/internal/pkg/server"
)

func TuneDownload(conn net.Conn) (res cmds.Result) {
	var targetMs int64 = 3000 // 3 seconds
	incBytes := 1048576       // 1 meg
	numBytes := 104857        // 0.1 megs
	maxBytes := 31457280      // 30 megs

	for {
		res = cmds.Download(conn, numBytes)
		fmt.Printf("Results: %#v\n", res)
		if res.DurationMs >= targetMs {
			break
		}
		if numBytes >= maxBytes {
			break
		}
		numBytes = numBytes + incBytes
	}
	return res
}

// GetAllServers returns all recommended servers
func GetAllServers() (servers []server.Server) {
	return server.GetAllServers()
}

// GetBestServer gets the first in the list
func GetBestServer() (bestserver server.Server) {
	fmt.Println("Finding best server")
	var bestspeed int64 = 999
	servers := GetAllServers()
	for s := range servers {
		c := Connect(servers[s].Host)
		res := PingTest(c, 3)
		if res < bestspeed {
			bestspeed = res
			bestserver = servers[s]
		}
	}

	return bestserver
}

// Connect returns the initial connection to the testing server
func Connect(server string) (conn net.Conn) {
	return cmds.Connect(server)
}

// Version returns the protocol version of speedtest binary protocol
func Version(conn net.Conn) (version string) {
	return cmds.Version(conn)
}

// DownloadTest runs numtests download tests for numbytes requested bytes
func DownloadTest(conn net.Conn, numbytes []int, numtests int) (results float64) {
	var acc float64

	fmt.Printf("Download test: ")
	for i := range numbytes {
		for j := 0; j < numtests; j++ {
			if !viper.GetBool("true") {
				fmt.Printf(".")
			}
			res := cmds.Download(conn, numbytes[i])
			mbps := CalcMbps(res.Start, res.Finish, res.Bytes)
			acc = acc + mbps
		}
	}

	results = acc / float64(numtests)
	fmt.Printf("\n")
	return results
}

// UploadTest runs numtests upload tests of numbytes random bytes
func UploadTest(conn net.Conn, numbytes []int, numtests int) (results float64) {
	var acc float64

	fmt.Printf("Upload test: ")
	for i := range numbytes {
		for j := 0; j < numtests; j++ {
			if !viper.GetBool("true") {
				fmt.Printf(".")
			}
			res := cmds.Upload(conn, numbytes[i])
			mbps := CalcMbps(res.Start, res.Finish, res.Bytes)
			acc = acc + mbps
		}
	}

	results = acc / float64(numtests)
	fmt.Printf("\n")
	return results
}

// PingTest gets roundtrip time to issue the "PING" command
func PingTest(conn net.Conn, numtests int) (results int64) {
	var acc int64

	fmt.Printf("Ping test: ")
	for i := 0; i < numtests; i++ {
		if !viper.GetBool("true") {
			fmt.Printf(".")
		}
		res := cmds.Ping(conn)
		acc = acc + res
	}

	results = acc / int64(numtests)
	fmt.Printf("\n")
	return results
}

func CalcMbps(start time.Time, finish time.Time, numbytes int) (mbps float64) {
	diff := finish.Sub(start)
	secs := float64(diff.Nanoseconds()) / float64(1000000000)
	megabits := float64(numbytes) / float64(125000)
	mbps = megabits / secs
	return mbps
}
