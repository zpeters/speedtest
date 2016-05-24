package tests

import (
	"fmt"
	"log"
	"strings"

	"github.com/zpeters/speedtest/internal/misc"
	"github.com/zpeters/speedtest/internal/print"
	"github.com/zpeters/speedtest/internal/sthttp"

	"github.com/spf13/viper"
)

// DownloadTest will perform the "normal" speedtest download test
func DownloadTest(server sthttp.Server) float64 {
	var urls []string
	var maxSpeed float64
	var avgSpeed float64

	// http://speedtest1.newbreakcommunications.net/speedtest/speedtest/
	for size := range viper.Get("dlsizes").([]int) {
		url := server.URL
		splits := strings.Split(url, "/")
		baseURL := strings.Join(splits[1:len(splits)-1], "/")
		randomImage := fmt.Sprintf("random%dx%d.jpg", viper.Get("dlsizes").([]int)[size], viper.Get("dlsizes").([]int)[size])
		downloadURL := "http:/" + baseURL + "/" + randomImage
		urls = append(urls, downloadURL)
	}

	if !viper.GetBool("quiet") && !viper.GetBool("report") {
		log.Printf("Testing download speed")
	}

	for u := range urls {

		if viper.GetBool("debug") {
			log.Printf("Download Test Run: %s\n", urls[u])
		}
		dlSpeed := sthttp.DownloadSpeed(urls[u])
		if !viper.GetBool("quiet") && !viper.GetBool("debug") && !viper.GetBool("report") {
			fmt.Printf(".")
		}
		if viper.GetBool("debug") {
			log.Printf("Dl Speed: %v\n", dlSpeed)
		}

		if viper.GetString("algotype") == "max" {
			if dlSpeed > maxSpeed {
				maxSpeed = dlSpeed
			}
		} else {
			avgSpeed = avgSpeed + dlSpeed
		}

	}

	if !viper.GetBool("quiet") && !viper.GetBool("report") {
		fmt.Printf("\n")
	}

	if viper.GetString("algotype") != "max" {
		return avgSpeed / float64(len(urls))
	}
	return maxSpeed

}

// UploadTest runs a "normal" speedtest upload test
func UploadTest(server sthttp.Server) float64 {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var maxSpeed float64
	var avgSpeed float64

	for size := range viper.Get("ulsizes").([]int) {
		ulsize = append(ulsize, viper.Get("ulsizes").([]int)[size])
	}

	if !viper.GetBool("quiet") && !viper.GetBool("report") {
		log.Printf("Testing upload speed")
	}

	for i := 0; i < len(ulsize); i++ {
		if viper.GetBool("debug") {
			log.Printf("Upload Test Run: %v\n", i)
		}
		r := misc.Urandom(ulsize[i])
		ulSpeed := sthttp.UploadSpeed(server.URL, "text/xml", r)
		if !viper.GetBool("quiet") && !viper.GetBool("debug") && !viper.GetBool("report") {
			fmt.Printf(".")
		}
		if viper.GetBool("debug") {
			log.Printf("Ul Amount: %v bytes\n", len(r))
			log.Printf("Ul Speed: %vMbps\n", ulSpeed)
		}

		if viper.GetString("algotype") == "max" {
			if ulSpeed > maxSpeed {
				maxSpeed = ulSpeed
			}
		} else {
			avgSpeed = avgSpeed + ulSpeed
		}

	}

	if !viper.GetBool("quiet") && !viper.GetBool("report") {
		fmt.Printf("\n")
	}

	if viper.GetString("algotype") != "max" {
		return avgSpeed / float64(len(viper.Get("ulsizes").([]int)))
	}
	return maxSpeed
}

// FindServer will find a specific server in the servers list
func FindServer(id string, serversList []sthttp.Server) sthttp.Server {
	var foundServer sthttp.Server
	for s := range serversList {
		if serversList[s].ID == id {
			foundServer = serversList[s]
		}
	}
	if foundServer.ID == "" {
		log.Fatalf("Cannot locate server Id '%s' in our list of speedtest servers!\n", id)
	}
	return foundServer
}

// ListServers prints a list of all "global" servers
func ListServers() {
	if viper.GetBool("debug") {
		fmt.Printf("Loading config from speedtest.net\n")
	}
	sthttp.CONFIG = sthttp.GetConfig()
	if viper.GetBool("debug") {
		fmt.Printf("\n")
	}

	if viper.GetBool("debug") {
		fmt.Printf("Getting servers list...")
	}
	allServers := sthttp.GetServers()
	if viper.GetBool("debug") {
		fmt.Printf("(%d) found\n", len(allServers))
	}
	for s := range allServers {
		server := allServers[s]
		print.Server(server)
	}
}
