package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
        "runtime"
)

import (
	"github.com/codegangsta/cli"
)

import (
	"github.com/zpeters/speedtest/debug"
	"github.com/zpeters/speedtest/misc"
	"github.com/zpeters/speedtest/sthttp"
)

var VERSION = "0.07.5"

var NUMCLOSEST int
var NUMLATENCYTESTS int
var TESTSERVERID = ""
var PINGONLY bool = false
var REPORTCHAR = ""
var ALGOTYPE = ""

func init_foo() {
	rand.Seed(time.Now().UTC().UnixNano())

	listFlag := flag.Bool("l", false, "\tList servers (hint use 'grep' or 'findstr' to locate a\n\t\t  server ID to use for '-s'")
	flag.BoolVar(&PINGONLY, "p", false, "\tPing only mode")
	flag.StringVar(&TESTSERVERID, "s", "", "\tSpecify a server ID to use")
	flag.StringVar(&ALGOTYPE, "a", "max", "\tSpecify the measurement method to use ('max', 'avg')")
	flag.IntVar(&NUMCLOSEST, "nc", 3, "\tNumber of geographically close servers to test to find\n\t\t  the optimal server")
	flag.IntVar(&NUMLATENCYTESTS, "nl", 5, "\tNumber of latency tests to perform to determine\n\t\t  which server is the fastest")
	verFlag := flag.Bool("v", false, "\tDisplay version")
	reportFlag := flag.Bool("r", false, "\tReporting mode output, minimal output with '|' for\n\t\t  separators, use '-rc' to change separator characters.\n\t\t  Reports the following: Server ID, Server Name (Location),\n\t\t  Ping time in ms, Download speed in kbps, Upload speed in kbps")
	flag.StringVar(&REPORTCHAR, "rc", "|", "\tCharacter to use to separate fields in report mode (-r)")

	flag.Parse()

	if *verFlag == true {
		fmt.Printf("%s - Version: %s\n", os.Args[0], VERSION)
	
		os.Exit(0)
	}

	if *reportFlag == true {
		debug.REPORT = true
		debug.QUIET = true
	}

	if debug.DEBUG {
		log.Printf("Debugging on...\n")
		debug.QUIET = false
	}

	if *listFlag == true {
		if debug.DEBUG {
			fmt.Printf("Loading config from speedtest.net\n")
		}
		sthttp.CONFIG = sthttp.GetConfig()
		if debug.DEBUG {
			fmt.Printf("\n")
		}

		if debug.DEBUG {
			fmt.Printf("Getting servers list...")
		}
		allServers := sthttp.GetServers()
		if debug.DEBUG {
			fmt.Printf("(%d) found\n", len(allServers))
		}
		for s := range allServers {
			server := allServers[s]
			printServer(server)
		}
		os.Exit(0)
	}
}

func downloadTest(server sthttp.Server) float64 {
	var urls []string
	var maxSpeed float64
	var avgSpeed float64

	// http://speedtest1.newbreakcommunications.net/speedtest/speedtest/
	dlsizes := []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

	for size := range dlsizes {
		url := server.Url
		splits := strings.Split(url, "/")
		baseUrl := strings.Join(splits[1:len(splits)-1], "/")
		randomImage := fmt.Sprintf("random%dx%d.jpg", dlsizes[size], dlsizes[size])
		downloadUrl := "http:/" + baseUrl + "/" + randomImage
		urls = append(urls, downloadUrl)
	}

	if !debug.QUIET {
		log.Printf("Testing download speed")
	}

	for u := range urls {

		if debug.DEBUG { fmt.Printf("Download Test Run: %s\n", urls[u])}
		dlSpeed := sthttp.DownloadSpeed(urls[u])
		if !debug.QUIET && !debug.DEBUG {
			fmt.Printf(".")
		}
		if debug.DEBUG {
			log.Printf("Dl Speed: %v\n", dlSpeed)
		}
		
		if ALGOTYPE == "max" {
			if dlSpeed > maxSpeed {
				maxSpeed = dlSpeed
			}
		} else {
			avgSpeed = avgSpeed + dlSpeed
		}
			
	}

	if !debug.QUIET {
		fmt.Printf("\n")
	}

	if ALGOTYPE == "max" {
		return maxSpeed
	} else {
		return avgSpeed / float64(len(urls))
	}
}

func uploadTest(server sthttp.Server) float64 {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var maxSpeed float64
	var avgSpeed float64
	
	ulsizesizes := []int{
		int(0.25 * 1024 * 1024),
		int(0.5 * 1024 * 1024),
		int(1.0 * 1024 * 1024),
		int(1.5 * 1024 * 1024),
		int(2.0 * 1024 * 1024),
	}

	for size := range ulsizesizes {
		ulsize = append(ulsize, ulsizesizes[size])
	}

	if !debug.QUIET { log.Printf("Testing upload speed") }
	
	for i:=0; i<len(ulsize); i++ {
		if debug.DEBUG { fmt.Printf("Upload Test Run: %v\n", i)}
		r := misc.Urandom(ulsize[i])
		ulSpeed := sthttp.UploadSpeed(server.Url, "text/xml", r)
		if !debug.QUIET && !debug.DEBUG {
			fmt.Printf(".")
		}


		if ALGOTYPE == "max" {
			if ulSpeed > maxSpeed {
				maxSpeed = ulSpeed
			}
		} else {
			avgSpeed = avgSpeed + ulSpeed
		}

	}

	if !debug.QUIET {
		fmt.Printf("\n")
	}



	if ALGOTYPE == "max" {
		return maxSpeed
	} else {
		return avgSpeed / float64(len(ulsizesizes))
	}
}

func findServer(id string, serversList []sthttp.Server) sthttp.Server {
	var foundServer sthttp.Server
	for s := range serversList {
		if serversList[s].Id == id {
			foundServer = serversList[s]
		}
	}
	if foundServer.Id == "" {
		log.Fatalf("Cannot locate server Id '%s' in our list of speedtest servers!\n", id)
	}
	return foundServer
}

func printServer(server sthttp.Server) {
	fmt.Printf("%-4s | %s (%s, %s)\n", server.Id, server.Sponsor, server.Name, server.Country)
}

func printServerReport(server sthttp.Server) {
	fmt.Printf("%s%s%s%s%s(%s,%s)%s", time.Now(), REPORTCHAR, server.Id, REPORTCHAR, server.Sponsor, server.Name, server.Country, REPORTCHAR)
}

func environmentReport() {
	fmt.Printf("Arch: %v\n", runtime.GOARCH) 
	fmt.Printf("OS: %v\n", runtime.GOOS) 
	fmt.Printf("IP: %v\n", sthttp.CONFIG.Ip) 
	fmt.Printf("Lat: %v\n", sthttp.CONFIG.Lat) 
	fmt.Printf("Lon: %v\n", sthttp.CONFIG.Lon) 
	fmt.Printf("ISP: %v\n", sthttp.CONFIG.Isp)
	fmt.Printf("-------------------------------\n")
	fmt.Printf("Debug: %v\n", debug.DEBUG)
	fmt.Printf("Quiet: %v\n", debug.QUIET)
}

func main() {
	// seeding randomness
	rand.Seed(time.Now().UTC().UnixNano())

	// setting up cli settings
	app := cli.NewApp()
	app.Name = "speedtest"
	app.Usage = "Unofficial command line interface to speedtest.net (https://github.com/zpeters/speedtest)"
	app.Author = "Zach Peters - zpeters@gmail.com - github.com/zpeters"
	app.Version = VERSION

	// setup cli flags
	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "debug, d",
			Usage: "Turn on debugging",
		},
		cli.BoolFlag{
			Name: "quiet, q",
			Usage: "Quiet mode",
		},
	}


	// setup the app aciton
	app.Action = func(c *cli.Context) {
		fmt.Printf("Bool: %s\n", c.Bool("debug"))
		fmt.Printf("Queit: %s\n", c.Bool("quiet"))
		if c.Bool("debug") {
			debug.DEBUG = true
		}
		if c.Bool("quiet") {
			debug.QUIET = true
		}
		
	}
	
	
	// var testServer sthttp.Server

	// if debug.DEBUG {
	// 	fmt.Printf("Loading config from speedtest.net\n")
	// }
	// sthttp.CONFIG = sthttp.GetConfig()

	// if debug.DEBUG { fmt.Printf("Environment report\n") }
	//app.Action = func(c *cli.Context) {
	//	if c.Bool("debug") {
	//		environmentReport()
//		}
//	}


	// run the app
	app.Run(os.Args)
	

	environmentReport()
	// if debug.DEBUG { fmt.Printf("Getting servers list...") }
	// allServers := sthttp.GetServers()
	// if debug.DEBUG {
	// 	fmt.Printf("(%d) found\n", len(allServers))
	// }

	
	// if TESTSERVERID != "" {
	// 	testServer = findServer(TESTSERVERID, allServers)

	// 	if !debug.REPORT {
	// 		printServer(testServer)
	// 	} else {
	// 		printServerReport(testServer)
	// 	}

	// 	if !debug.QUIET && !debug.REPORT {
	// 		fmt.Printf("Testing latency...\n")
	// 	}
	// 	testServer.Latency = sthttp.GetLatency(testServer, NUMLATENCYTESTS, ALGOTYPE)
	// } else {
	// 	closestServers := sthttp.GetClosestServers(allServers)
	// 	if !debug.QUIET && !debug.REPORT {
	// 		log.Printf("Finding fastest server..\n")
	// 	}
	// 	testServer = sthttp.GetFastestServer(NUMCLOSEST, NUMLATENCYTESTS, closestServers, ALGOTYPE)

	// 	if !debug.REPORT {
	// 		printServer(testServer)
	// 	} else {
	// 		printServerReport(testServer)
	// 	}

	// 	if debug.DEBUG {
	// 		fmt.Printf("\n")
	// 	}
	// }

	// if PINGONLY {		
	// 	if !debug.REPORT {
	// 		if ALGOTYPE == "max" {
	// 			fmt.Printf("Ping (Lowest): %3.2f ms\n", testServer.Latency)
	// 		} else {
	// 			fmt.Printf("Ping (Avg): %3.2f ms\n", testServer.Latency)
	// 		}
	// 	} else {
	// 		fmt.Printf("%3.2f (Max)\n", testServer.Latency)
	// 	}
	// } else {
	// 	dmbps := downloadTest(testServer)
	// 	umbps := uploadTest(testServer)
	// 	if !debug.REPORT {
	// 		if ALGOTYPE == "max" {
	// 			fmt.Printf("Ping (Lowest): %3.2f ms | Download (Max): %3.2f Mbps | Upload (Max): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
	// 		} else {
	// 			fmt.Printf("Ping (Avg): %3.2f ms | Download (Avg): %3.2f Mbps | Upload (Avg): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
	// 		}
	// 	} else {
	// 		dkbps := dmbps * 1000
	// 		ukbps := umbps * 1000
	// 		fmt.Printf("%3.2f%s%d%s%d\n", testServer.Latency, REPORTCHAR, int(dkbps), REPORTCHAR, int(ukbps))
	// 		fmt.Printf("%3.2f%s%d%s%d\n", testServer.Latency, REPORTCHAR, int(dkbps), REPORTCHAR, int(ukbps))
	// 	}
	// }
}
