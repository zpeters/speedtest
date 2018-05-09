/*
speedtest is an unofficial commandline interface to speedtest.net

Version 1.0 was designed as an "app only".  Version 2.0 will make a cleaner split between libraries and interface
*/

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/zpeters/speedtest/internal/print"
	"github.com/zpeters/speedtest/internal/speedtests"
	"github.com/zpeters/speedtest/internal/sthttp"

	"github.com/dchest/uniuri"
	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Version placeholder, injected in Makefile
var Version string

func runTest(c *cli.Context, stClient *sthttp.Client, tester *speedtests.Tester) {
	// create our server object and load initial config
	var testServer sthttp.Server

	config, err := stClient.GetConfig()
	if err != nil {
		log.Printf("Cannot get speedtest config\n")
		log.Fatal(err)
	}
	stClient.Config = &config

	// if we are *not* running a report then say hello to everyone
	if !tester.Report {
		fmt.Printf("github.com/zpeters/speedtest -- unofficial cli for speedtest.net\n")
	}

	// if we are in debug mode print outa an environment report
	if stClient.Debug {
		print.EnvironmentReport(stClient)
	}

	// get all possible servers (excluding blacklisted)
	if stClient.Debug {
		log.Printf("Getting all servers for our test list")
	}
	var allServers []sthttp.Server
	if c.String("mini") == "" {
		allServers, err = stClient.GetServers()
		if err != nil {
			log.Fatal(err)
		}
	}

	// if a mini speedtest installation was specified, use that...
	if c.String("mini") != "" {

		//construct testserver object manually
		u, err := url.Parse(c.String("mini"))
		if err != nil {
			log.Fatalf("Speedtest mini server URL is not a valid URL: %s", err)
		}

		if stClient.Debug {
			log.Printf("Using Mini Server '%s'", c.String("mini"))
		}
		testServer.URL = c.String("mini")
		if !strings.HasSuffix(c.String("mini"), "/") {
			testServer.URL += "/"
		}
		testServer.URL += "speedtest/upload.php"
		testServer.Name = u.Host
		testServer.Sponsor = "speedtest-mini"
		testServer.ID = "0"

		testServer.Latency, err = stClient.GetLatency(testServer, stClient.GetLatencyURL(testServer))
		if err != nil {
			log.Fatal(err)
		}

		// if they specified a specific speedtest.net server, test against that...
	} else if c.String("server") != "" {
		if stClient.Debug {
			log.Printf("Server '%s' specified, getting info...", c.String("server"))
		}
		// find server and load latency report
		testServer = tester.FindServer(c.String("server"), allServers)
		// load latency
		testServer.Latency, err = stClient.GetLatency(testServer, stClient.GetLatencyURL(testServer))
		if err != nil {
			log.Fatal(err)
		}

		if !tester.Report {
			fmt.Printf("Server: %s - %s (%s)\n", testServer.ID, testServer.Name, testServer.Sponsor)
		}

		// ...otherwise get a list of all servers sorted by distance...
	} else {
		if stClient.Debug {
			log.Printf("Getting closest servers...")
		}
		closestServers := stClient.GetClosestServers(allServers)
		if stClient.Debug {
			log.Printf("Getting the fastests of our closest servers...")
		}
		// ... and get the fastests NUMCLOSEST ones
		testServer = stClient.GetFastestServer(closestServers)
		if !viper.GetBool("report") {
			fmt.Printf("Server: %s - %s (%s)\n", testServer.ID, testServer.Name, testServer.Sponsor)
		}
	}

	// if ping only then just output latency results and exit nicely...
	if c.Bool("ping") {
		if c.Bool("report") {
			if viper.GetString("algotype") == "max" {
				fmt.Printf("%3.2f (Lowest)\n", testServer.Latency)
			} else {
				fmt.Printf("%3.2f (Avg)\n", testServer.Latency)
			}
		} else {
			if viper.GetString("algotype") == "max" {
				fmt.Printf("Ping (Lowest): %3.2f ms\n", testServer.Latency)
			} else {
				fmt.Printf("Ping (Avg): %3.2f ms\n", testServer.Latency)
			}
		}
		os.Exit(0)
		// ...otherwise run our full test
	} else {
		var dmbps float64
		var umbps float64

		if !viper.GetBool("report") {
			if c.Bool("downloadonly") {
				dmbps = tester.Download(testServer)
			} else if c.Bool("uploadonly") {
				umbps = tester.Upload(testServer)
			} else {
				dmbps = tester.Download(testServer)
				umbps = tester.Upload(testServer)
			}
			if viper.GetString("algotype") == "max" {
				fmt.Printf("Ping (Lowest): %3.2f ms | Download (Max): %3.2f Mbps | Upload (Max): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
			} else {
				fmt.Printf("Ping (Avg): %3.2f ms | Download (Avg): %3.2f Mbps | Upload (Avg): %3.2f Mbps\n", testServer.Latency, dmbps, umbps)
			}

		} else {

			fmt.Printf("%s%s%s%s\"%s (%s, %s)\"%s", time.Now().Format("2006-01-02 15:04:05 -0700"), viper.GetString("reportchar"), testServer.ID, viper.GetString("reportchar"), testServer.Sponsor, testServer.Name, testServer.Country, viper.GetString("reportchar"))
			fmt.Printf("%3.2f%s", testServer.Latency, viper.GetString("reportchar"))

			if c.Bool("downloadonly") {
				dmbps = tester.Download(testServer)
				dkbps := dmbps * 1000
				fmt.Printf("%d\n", int(dkbps))
			} else if c.Bool("uploadonly") {
				umbps = tester.Upload(testServer)
				ukbps := umbps * 1000
				fmt.Printf("%d\n", int(ukbps))
			} else {
				dmbps = tester.Download(testServer)
				dkbps := dmbps * 1000
				fmt.Printf("%d%s", int(dkbps), viper.GetString("reportchar"))

				umbps = tester.Upload(testServer)
				ukbps := umbps * 1000
				fmt.Printf("%d\n", int(ukbps))
			}
		}
	}
}

func init() {
	viper.SetDefault("debug", false)
	viper.SetDefault("quiet", false)
	viper.SetDefault("report", false)
	viper.SetDefault("numclosest", 3)
	viper.SetDefault("numlatencytests", 5)
	viper.SetDefault("reportchar", "|")
	viper.SetDefault("algotype", "max")
	viper.SetDefault("httptimeout", 15)
	viper.SetDefault("dlsizes", []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000})
	viper.SetDefault("ulsizes", []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)})
	viper.SetDefault("speedtestconfigurl", "http://c.speedtest.net/speedtest-config.php?x="+uniuri.New())
	viper.SetDefault("speedtestserversurl", "http://c.speedtest.net/speedtest-servers-static.php?x="+uniuri.New())
	viper.SetDefault("useragent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.21 Safari/537.36")
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
	app.Version = Version

	// setup cli flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "algo, a",
			Usage: "Specify the measurement method to use ('max', 'avg')",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Turn on debugging",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "List available servers",
		},
		cli.BoolFlag{
			Name:  "update, u",
			Usage: "Check for a new version of speedtest",
		},
		cli.BoolFlag{
			Name:  "ping, p",
			Usage: "Ping only mode",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Quiet mode",
		},
		cli.BoolFlag{
			Name:  "report, r",
			Usage: "Reporting mode output, minimal output with '|' for separators, use '--rc'\n\t\tto change separator characters. Reports the following: Server ID, \n\t\tServer Name (Location), Ping time in ms, Download speed in kbps, Upload speed in kbps",
		},
		cli.BoolFlag{
			Name:  "downloadonly, do",
			Usage: "Only perform download test",
		},
		cli.BoolFlag{
			Name:  "uploadonly, uo",
			Usage: "Only perform upload test",
		},
		cli.StringFlag{
			Name:  "reportchar, rc",
			Usage: "Set the report separator. Example: --rc=','",
		},
		cli.StringFlag{
			Name:  "server, s",
			Usage: "Use a specific server",
		},
		cli.StringSliceFlag{
			Name:  "blacklist, b",
			Usage: "Blacklist a server.  Use this multiple times for more than one server",
		},
		cli.StringFlag{
			Name:  "mini, m",
			Usage: "URL of speedtest mini server",
		},
		cli.StringFlag{
			Name:  "useragent, ua",
			Usage: "Specify a useragent string",
		},
		cli.IntFlag{
			Name:  "numclosest, nc",
			Value: viper.GetInt("numclosest"),
			Usage: "Number of 'closest' servers to find",
		},
		cli.IntFlag{
			Name:  "httptimeout, t",
			Value: viper.GetInt("httptimeout"),
			Usage: "Timeout (seconds) for http connections",
		},
		cli.IntFlag{
			Name:  "numlatency, nl",
			Value: viper.GetInt("numlatencytests"),
			Usage: "Number of latency tests to perform",
		},
		cli.StringFlag{
			Name:  "interface, I",
			Usage: "Source IP address or name of an interface",
		},
	}

	// toggle our switches and setup variables
	app.Action = func(c *cli.Context) {
		// just check the version if that is what they want
		if c.Bool("update") {
			// Check if there is an update
			client := github.NewClient(nil)
			ctx := context.Background()
			latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, "zpeters", "speedtest")
			if err != nil {
				log.Fatalf("github call: %s", err)
			}
			githubTag := *latestRelease.TagName
			if Version != githubTag {
				fmt.Printf("New version %s available at https://github.com/zpeters/speedtest/releases\n", githubTag)
			} else {
				fmt.Printf("You are up to date\n")
			}
			os.Exit(0)
		}
		// set our flags
		if c.Bool("debug") {
			viper.Set("debug", true)
		}
		if c.Bool("quiet") {
			viper.Set("quiet", true)
		}
		if c.Bool("report") {
			viper.Set("report", true)
		}
		if c.String("algo") != "" {
			if c.String("algo") == "max" {
				viper.Set("algotype", "max")
			} else if c.String("algo") == "avg" {
				viper.Set("algotype", "avg")
			} else {
				fmt.Printf("** Invalid algorithm '%s'\n", c.String("algo"))
				os.Exit(1)
			}
		}
		viper.Set("numclosest", c.Int("numclosest"))
		viper.Set("numlatencytests", c.Int("numlatency"))
		viper.Set("httptimeout", c.Int("httptimeout"))
		if c.String("reportchar") != "" {
			viper.Set("reportchar", c.String("reportchar"))
		}
		if c.String("interface") != "" {
			viper.Set("interface", c.String("interface"))
		}
		if len(c.StringSlice("blacklist")) > 0 {
			viper.Set("blacklist", c.StringSlice("blacklist"))
		}

		stClient := sthttp.NewClient(
			&sthttp.SpeedtestConfig{
				ConfigURL:       viper.GetString("speedtestconfigurl"),
				ServersURL:      viper.GetString("speedtestserversurl"),
				AlgoType:        viper.GetString("algotype"),
				NumClosest:      viper.GetInt("numclosest"),
				NumLatencyTests: viper.GetInt("numlatencytests"),
				Interface:       viper.GetString("interface"),
				Blacklist:       viper.GetStringSlice("blacklist"),
				UserAgent:       viper.GetString("useragent"),
			},
			&sthttp.HTTPConfig{
				HTTPTimeout: viper.GetDuration("httptimeout") * time.Second,
			},
			viper.GetBool("debug"),
			viper.GetString("reportchar"))

		tester := speedtests.NewTester(
			stClient,
			viper.Get("dlsizes").([]int),
			viper.Get("ulsizes").([]int),
			viper.GetBool("quiet"),
			viper.GetBool("report"))

		// run a oneshot list
		if c.Bool("list") {
			tester.ListServers(stClient.SpeedtestConfig.ConfigURL, stClient.SpeedtestConfig.ServersURL, stClient.SpeedtestConfig.Blacklist)
			os.Exit(0)
		}

		// run our test
		runTest(c, stClient, tester)

		// exit nicely
		os.Exit(0)
	}

	// run the app
	app.Run(os.Args)
}
