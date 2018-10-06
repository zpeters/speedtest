package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zpeters/speedtest/internal/app"
	"github.com/zpeters/speedtest/internal/pkg/cmds"
)

func init() {
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
}

func main() {
	var seedbytes int = 1000000
	var numping int = 10
	var timelimit int = 7

	start := time.Now()

	server, err := app.GetBestServer()
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"package":  "main",
			"function": "main",
		}).Fatal()
	}
	fmt.Printf("Found best server: (%s) %s - %s\n", server.ID, server.Name, server.Sponsor)
	conn := cmds.Connect(server.Host)
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"package":  "main",
			"function": "main",
		}).Fatal()
	}

	fmt.Printf("Speedtest protocol version: %s\n", cmds.Version(conn))

	ping := app.PingTest(conn, numping)

	log.WithFields(log.Fields{
		"seed":     seedbytes,
		"package":  "main",
		"funciton": "main",
	}).Debug("DownloadBytes")
	download := app.DownloadTest(conn, seedbytes, timelimit)

	log.WithFields(log.Fields{
		"seed":     seedbytes,
		"package":  "main",
		"funciton": "main",
	}).Debug("UploadBytes")
	upload := app.UploadTest(conn, seedbytes, timelimit)

	complete := time.Now()
	elapsed := complete.Sub(start)

	fmt.Printf("Ping results: %d ms\n", ping)
	fmt.Printf("Download results: %f mbps\n", download)
	fmt.Printf("Upload results: %f mbps\n", upload)
	fmt.Printf("Took: %s\n", elapsed)
}
