package main

import (
	"fmt"
)

import (
	log "github.com/sirupsen/logrus"
)

import (
	"github.com/zpeters/speedtest/internal/app"
	"github.com/zpeters/speedtest/internal/pkg/cmds"
)

func init() {
	log.SetLevel(log.WarnLevel)
	//log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
}

func main() {

	server, err := app.GetBestServer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"package": "main",
			"function": "main",
		}).Fatal()
	}
	fmt.Printf("Found best server: (%s) %s - %s\n", server.ID, server.Name, server.Sponsor)
	conn := cmds.Connect(server.Host)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"package": "main",
			"function": "main",
		}).Fatal()
	}

	fmt.Printf("Speedtest protocol version: %s\n", cmds.Version(conn))

	ping := app.PingTest(conn, 20)

	downloadBytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000, 2000000}
	log.WithFields(log.Fields{
		"testrange": downloadBytes,
		"package": "main",
		"funciton": "main",
	}).Debug("DownloadBytes")
	download := app.DownloadTest(conn, downloadBytes, 4)

	uploadBytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000}
	log.WithFields(log.Fields{
		"testrange": uploadBytes,
		"package": "main",
		"funciton": "main",
	}).Debug("UploadBytes")
	upload := app.UploadTest(conn, uploadBytes, 4)

	fmt.Printf("Ping results: %d ms\n", ping)
	fmt.Printf("Download results: %f mbps\n", download)
	fmt.Printf("Upload results: %f mbps\n", upload)
}
