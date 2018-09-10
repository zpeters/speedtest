package app

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/viper"
	"github.com/zpeters/speedtest/internal/pkg/cmds"
	"github.com/zpeters/speedtest/internal/pkg/server"
)

// GetBestServer gets the first in the list
func GetBestServer() (bestserver server.Server, err error) {
	fmt.Println("Finding best server")
	var bestspeed int64 = 999
	servers, err := server.GetAllServers()
	if err != nil {
		log.Fatal(err)
	}
	for s := range servers {
		c := cmds.Connect(servers[s].Host)
		res := PingTest(c, 3)
		servers[s].BestTestPing = res
		if res < bestspeed {
			bestspeed = res
			bestserver = servers[s]
		}
	}

	return bestserver, err
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
