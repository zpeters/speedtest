package main
import (
	"fmt"
	"log"
	"time"
	"os"
	"strings"
	"math/rand"
	"flag"
)

import (
	"speedtest/misc"
	"speedtest/debug"
	"speedtest/sthttp"
)

var NUMCLOSEST = 5
var NUMLATENCYTESTS = 5
var VERSION = "0.02"



func init() {
	flag.BoolVar(&debug.DEBUG, "d", false, "Turn on debugging")
	verFlag := flag.Bool("v", false, "Display version")
	flag.Parse()
	if *verFlag == true {
		fmt.Printf("%s - Version: %s\n", os.Args[0], VERSION)
		os.Exit(0)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	if debug.DEBUG { log.Printf("Debugging on...\n") }
}

func downloadTest(server sthttp.Server) float64 {
	var urls []string
	var speedAcc float64
	var numTests = 4
	dlsizes := []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	//dlsizes := []int{350, 500, 750}



	// generate the size urls
	for size := range dlsizes {
		for i := 0; i<numTests; i++ {
			url := server.Url
			splits := strings.Split(url, "/")
			baseUrl := strings.Join(splits[1:len(splits) -1], "/")
			randomImage := fmt.Sprintf("random%dx%d.jpg", dlsizes[size], dlsizes[size])
			downloadUrl := "http:/" + baseUrl + "/" + randomImage
			urls = append(urls, downloadUrl)
		}
	}	


	fmt.Printf("\tRunning %d tests, %d megs total\n", numTests, len(urls))

	// test the urls
	for u := range urls {
		if debug.DEBUG { fmt.Printf("Download test %d\n", u) }
		dlSpeed := sthttp.DownloadSpeed(urls[u])
		if debug.DEBUG { fmt.Printf("\tDownload speed: %f Mbps\n", dlSpeed) }
		speedAcc = speedAcc + dlSpeed
	}
	

	mbps := (speedAcc / float64(len(urls)))
	return mbps
}


func uploadTest(server sthttp.Server) float64 {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var ulSpeedAcc float64

	//ulsizesizes := []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024)}
	ulsizesizes := []int{int(0.25 * 1024 * 1024)}
	
	for size := range ulsizesizes {
		for i := 0; i<25; i++ {
			ulsize = append(ulsize, ulsizesizes[size])
		}
	}

	fmt.Printf("\tRunning %d tests - %d Megs total\n", len(ulsize), len(ulsize))
	

	for i:=0; i<len(ulsize); i++ {
		if debug.DEBUG { fmt.Printf("Ulsize: %d\n", ulsize[i]) }
		r := misc.Urandom(ulsize[i])
		ulSpeed := sthttp.UploadSpeed(server.Url, "text/xml", r)
		if debug.DEBUG { fmt.Printf("\tUpload speed: %f Mbps\n", ulSpeed) }
		ulSpeedAcc = ulSpeedAcc + ulSpeed
	}
	
	mbps := ulSpeedAcc / float64(len(ulsize))
	return mbps
}


func main() {
	fmt.Printf("Loading config...\n")
	sthttp.CONFIG = sthttp.GetConfig()

	fmt.Printf("Getting servers list...")
	allServers := sthttp.GetServers()
	fmt.Printf("(%d) found\n", len(allServers))
	
	fmt.Printf("Finding %d closest servers...\n", NUMCLOSEST)
	// add an option for num closest?
	closestServers := sthttp.GetClosestServers(NUMCLOSEST, allServers)
	
	fmt.Printf("Finding fastest server - testing latency %d times...\n", NUMLATENCYTESTS)
	// add an option for num rums, test how many are necessary
	fastestServer := sthttp.GetFastestServer(NUMLATENCYTESTS, closestServers)
	fmt.Printf("Fastest Server: %s (%s - %s) - %s ping \n", fastestServer.Sponsor, fastestServer.Name, fastestServer.Country, fastestServer.AvgLatency)
	
	fmt.Printf("Starting download test...\n")
	dmbps := downloadTest(fastestServer)
	fmt.Printf("Average Download Speed: %f Mbps\n", dmbps)
	
	fmt.Printf("Starting Upload test...\n")
	umbps := uploadTest(fastestServer)
	fmt.Printf("Average Upload Speed: %f Mbps\n", umbps)
 	
}
