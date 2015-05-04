package main

import (
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

var NUMCLOSEST int = 3
var NUMLATENCYTESTS int = 5
var REPORTCHAR = "|"
var ALGOTYPE = "max"

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

func environmentReport(c *cli.Context) {
	log.Printf("Env Report")
	log.Printf("-------------------------------\n")
	log.Printf("[User Environment]\n")
	log.Printf("Arch: %v\n", runtime.GOARCH) 
	log.Printf("OS: %v\n", runtime.GOOS) 
	log.Printf("IP: %v\n", sthttp.CONFIG.Ip) 
	log.Printf("Lat: %v\n", sthttp.CONFIG.Lat) 
	log.Printf("Lon: %v\n", sthttp.CONFIG.Lon) 
	log.Printf("ISP: %v\n", sthttp.CONFIG.Isp)
	log.Printf("-------------------------------\n")
	log.Printf("[Settings]\n")
	if c.Bool("debug") {
		log.Printf("Debug (user): %v\n", debug.DEBUG)
	} else {
		log.Printf("Debug (default): %v\n", debug.DEBUG)
	}
	if c.Bool("quiet") {
		log.Printf("Quiet (user): %v\n", debug.QUIET)
	} else {
		log.Printf("Quiet (default): %v\n", debug.QUIET)
	}
	if c.Int("numclosest") == 0 {
		log.Printf("NUMCLOSEST (default): %v\n", NUMCLOSEST)
	} else {
		log.Printf("NUMCLOSEST (user): %v\n", NUMCLOSEST)

	}
	if c.Int("numlatency") == 0 {
		log.Printf("NUMLATENCYTESTS (default): %v\n", NUMLATENCYTESTS)
	} else {
		log.Printf("NUMLATENCYTESTS (user): %v\n", NUMLATENCYTESTS)
	}
	if c.String("server") == "" {
		log.Printf("server (default none specified)\n")
	} else {
		log.Printf("server (user): %s\n", c.String("server"))
	}
	if c.String("reportchar") == "" {
		log.Printf("reportchar (default): %s\n", REPORTCHAR)
	} else {
		log.Printf("reportchar (user): %s\n", c.String("reportchar"))
	}
	if c.String("algo") == "" {
		log.Printf("algo (default): %s\n", ALGOTYPE)
	} else {
		log.Printf("algo (user): %s\n", c.String("algo"))
	}
	log.Printf("--------------------------------\n")
	log.Printf("[Mode]\n")
	log.Printf("Report: %v\n", c.Bool("report"))
	log.Printf("List: %v\n", c.Bool("list"))
	log.Printf("Ping: %v\n", c.Bool("Ping"))
	log.Printf("-------------------------------\n")	

}

func listServers() {
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
}

func runTest(c *cli.Context) {
	// create our server object and load initial config
	var testServer sthttp.Server
	sthttp.CONFIG = sthttp.GetConfig()

	if debug.DEBUG {
		environmentReport(c)
	}

	// get all possible servers
	if debug.DEBUG {
		log.Printf("Getting all servers for our test list")
	}
	allServers := sthttp.GetServers()
	
	// if they specified a specific server, test against that...
	if c.String("server") != "" {
		if debug.DEBUG {
			log.Printf("Server '%s' specified, getting info...", c.String("server"))
		}
		// find server and load latency report
		testServer = findServer(c.String("server"), allServers)
		// load latency
		testServer.Latency = sthttp.GetLatency(testServer, NUMLATENCYTESTS, ALGOTYPE )

		fmt.Printf("Selected server: %s\n", testServer)
	// ...otherwise get a list of all servers sorted by distance...
	} else {
		if debug.DEBUG {
			log.Printf("Getting closest servers...")
		}
		closestServers := sthttp.GetClosestServers(allServers)
		if debug.DEBUG {
			log.Printf("Getting the fastests of our closest servers...")
		}
		// ... and get the fastests NUMCLOSEST ones
		testServer = sthttp.GetFastestServer(NUMCLOSEST, NUMLATENCYTESTS, closestServers, ALGOTYPE)
	}


	// Start printing our report
	if !debug.REPORT {
		printServer(testServer)
	} else {
		printServerReport(testServer)
	}

	// if ping only then just output latency results and exit nicely...
	if c.Bool("ping") {
		if c.Bool("report") {
			if ALGOTYPE == "max" {
				fmt.Printf("%3.2f (Lowest)\n", testServer.Latency)
			} else {
				fmt.Printf("%3.2f (Avg)\n", testServer.Latency)
			}
		} else {
			if ALGOTYPE == "max" {
				fmt.Printf("Ping (Lowest): %3.2f ms\n", testServer.Latency)
			} else {
				fmt.Printf("Ping (Avg): %3.2f ms\n", testServer.Latency)
			}
		}
		os.Exit(0)
	// ...otherwise run our full test
	} else {
		
		dmbps := downloadTest(testServer)
		umbps := uploadTest(testServer)
		if !debug.REPORT {
			if ALGOTYPE == "max" {
				fmt.Printf("Ping (Lowest): %3.2f ms | Download (Max): %3.2f Mbps | Upload (Max): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
			} else {
				fmt.Printf("Ping (Avg): %3.2f ms | Download (Avg): %3.2f Mbps | Upload (Avg): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
			}
		} else {
			dkbps := dmbps * 1000
			ukbps := umbps * 1000
			fmt.Printf("%3.2f%s%d%s%d\n", testServer.Latency, REPORTCHAR, int(dkbps), REPORTCHAR, int(ukbps))
		}
	}

}

func main() {
	// seeding randomness
	rand.Seed(time.Now().UTC().UnixNano())

	// set logging to stdout for global logger
	log.SetOutput(os.Stdout)

	// setting up cli settings
	app := cli.NewApp()
	app.Name = "speedtest"
	app.Usage = "Unofficial command line interface to speedtest.net (https://github.com/zpeters/speedtest)"
	app.Author = "Zach Peters - zpeters@gmail.com - github.com/zpeters"
	app.Version = VERSION

	// setup cli flags
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "algo, a",
			Usage: "Specify the measurement method to use ('max', 'avg')",
		},
		cli.BoolFlag{
			Name: "debug, d",
			Usage: "Turn on debugging",
		},
		cli.BoolFlag{
			Name: "list, l",
			Usage: "List available servers",
		},
		cli.BoolFlag{
			Name: "ping, p",
			Usage: "Ping only mode",
		},
		cli.BoolFlag{
			Name: "quiet, q",
			Usage: "Quiet mode",
		},
		cli.BoolFlag{
			Name: "report, r",
			Usage: "Reporting mode output, minimal output with '|' for separators, use '-rc' to change separator characters. Reports the following: Server ID, Server Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps",
		},
		cli.StringFlag{
			Name: "reportchar, rc",
			Usage: "Set the report separator",
		},
		cli.StringFlag{
			Name: "server, s",
			Usage: "Use a specific server",
		},
		cli.IntFlag{
			Name: "numclosest, nc",
			Value: NUMCLOSEST,
			Usage: "Number of 'closest' servers to find",
		},
		cli.IntFlag{
			Name: "numlatency, nl",
			Value: NUMLATENCYTESTS,
			Usage: "Number of latency tests to perform",
		},
			
	}

	// toggle our switches and setup variables
	app.Action = func(c *cli.Context) {
		// set our flags
		if c.Bool("debug") {
			debug.DEBUG = true
		}
		if c.Bool("quiet") {
			debug.QUIET = true
		}
		if c.Bool("report") {
			debug.REPORT = true
		}
		if c.String("algo") != "" {
			ALGOTYPE = c.String("algo")
		}
		NUMCLOSEST = c.Int("numclosest")
		NUMLATENCYTESTS = c.Int("numlatency")
		if c.String("reportchar") != "" {
			REPORTCHAR = c.String("reportchar")
		}

		// run a oneshot list
		if c.Bool("list") {
			listServers()
			os.Exit(0)
		}

		// run our test
		runTest(c)
	}
	// run the app
	app.Run(os.Args)
}
