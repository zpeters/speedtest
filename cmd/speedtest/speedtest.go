package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/zpeters/speedtest/internal/app"
)

func config() {
	viper.SetDefault("Debug", false)
	viper.SetDefault("DbLocation", "data.sqlite")
}

func main() {
	config()

	server := app.GetBestServer()
	fmt.Printf("Found best server: (%s) %s - %s\n", server.ID, server.Name, server.Sponsor)
	conn := app.Connect(server.Host)

	//log.Printf("Begin tuning...")
	//res := app.TuneDownload(conn)
	//fmt.Printf("Tuned Download: %#v", res)
	//mbps := app.CalcMbps(res.Start, res.Finish, res.Bytes )
	//log.Printf("MBPS: %#v\n", mbps)
	//log.Printf("Tuning complete...")

	fmt.Printf("Speedtest protocol version: %s\n", app.Version(conn))

	ping := app.PingTest(conn, 20)

	downloadBytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000, 2000000, 5000000, 7000000}
	if viper.GetBool("Debug") {
		log.Printf("downloadBytes: %#v\n", downloadBytes)
	}
	download := app.DownloadTest(conn, downloadBytes, 4)

	uploadBytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000, 2000000, 5000000, 7000000}
	if viper.GetBool("Debug") {
		log.Printf("uploadBytes: %#v\n", uploadBytes)
	}
	upload := app.UploadTest(conn, uploadBytes, 4)

	fmt.Printf("--| Results |---\n")
	fmt.Printf("Ping results: %d ms\n", ping)
	fmt.Printf("Download results: %f mbps\n", download)
	fmt.Printf("Upload results: %f mbps\n", upload)
}
