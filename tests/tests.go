package tests

import (
	"fmt"
	"log"
	"strings"

	"github.com/zpeters/speedtest/misc"
	"github.com/zpeters/speedtest/print"
	"github.com/zpeters/speedtest/sthttp"
)

var (
	// DefaultDLSizes defines the default download sizes
	DefaultDLSizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	// DefaultULSizes defines the default upload sizes
	DefaultULSizes = []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024), int(1.0 * 1024 * 1024), int(1.5 * 1024 * 1024), int(2.0 * 1024 * 1024)}
)

// Tester defines a Speedtester client tester
type Tester struct {
	Client   *sthttp.Client
	DLSizes  []int
	ULSizes  []int
	Quiet    bool
	Report   bool
	Debug    bool
	AlgoType string
}

func NewTester(client *sthttp.Client, dlsizes []int, ulsizes []int, quiet bool, report bool) *Tester {
	return &Tester{
		Client:  client,
		DLSizes: dlsizes,
		ULSizes: ulsizes,
		Quiet:   quiet,
		Report:  report,
	}
}

// Download will perform the "normal" speedtest download test
func (tester *Tester) Download(server sthttp.Server) float64 {
	var urls []string
	var maxSpeed float64
	var avgSpeed float64

	// http://speedtest1.newbreakcommunications.net/speedtest/speedtest/
	for size := range tester.DLSizes {
		url := server.URL
		splits := strings.Split(url, "/")
		baseURL := strings.Join(splits[1:len(splits)-1], "/")
		randomImage := fmt.Sprintf("random%dx%d.jpg", tester.DLSizes[size], tester.DLSizes[size])
		downloadURL := "http:/" + baseURL + "/" + randomImage
		urls = append(urls, downloadURL)
	}

	if !tester.Quiet && !tester.Report {
		log.Printf("Testing download speed")
	}

	for u := range urls {
		if tester.Debug {
			log.Printf("Download Test Run: %s\n", urls[u])
		}
		dlSpeed, err := tester.Client.DownloadSpeed(urls[u])
		if err != nil {
			log.Printf("Can't get download speed")
			log.Fatal(err)
		}
		if !tester.Quiet && !tester.Debug && !tester.Report {
			fmt.Printf(".")
		}
		if tester.Debug {
			log.Printf("Dl Speed: %v\n", dlSpeed)
		}

		if tester.AlgoType == "max" {
			if dlSpeed > maxSpeed {
				maxSpeed = dlSpeed
			}
		} else {
			avgSpeed = avgSpeed + dlSpeed
		}
	}

	if !tester.Quiet && !tester.Report {
		fmt.Printf("\n")
	}

	if tester.AlgoType != "max" {
		return avgSpeed / float64(len(urls))
	}
	return maxSpeed

}

// Upload runs a "normal" speedtest upload test
func (tester *Tester) Upload(server sthttp.Server) float64 {
	// https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	var ulsize []int
	var maxSpeed float64
	var avgSpeed float64

	for size := range tester.ULSizes {
		ulsize = append(ulsize, tester.ULSizes[size])
	}

	if !tester.Quiet && !tester.Report {
		log.Printf("Testing upload speed")
	}

	for i := 0; i < len(ulsize); i++ {
		if tester.Debug {
			log.Printf("Upload Test Run: %v\n", i)
		}
		r := misc.Urandom(ulsize[i])
		ulSpeed, err := tester.Client.UploadSpeed(server.URL, "text/xml", r)
		if err != nil {
			log.Fatal(err)
		}
		if !tester.Quiet && !tester.Debug && !tester.Report {
			fmt.Printf(".")
		}
		if tester.Debug {
			log.Printf("Ul Amount: %v bytes\n", len(r))
			log.Printf("Ul Speed: %vMbps\n", ulSpeed)
		}

		if tester.AlgoType == "max" {
			if ulSpeed > maxSpeed {
				maxSpeed = ulSpeed
			}
		} else {
			avgSpeed = avgSpeed + ulSpeed
		}

	}

	if !tester.Quiet && !tester.Report {
		fmt.Printf("\n")
	}

	if tester.AlgoType != "max" {
		return avgSpeed / float64(len(ulsize))
	}
	return maxSpeed
}

// FindServer will find a specific server in the servers list
func (tester *Tester) FindServer(id string, serversList []sthttp.Server) sthttp.Server {
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
func (tester *Tester) ListServers(configURL string, serversURL string, blacklist []string) (err error) {
	if tester.Debug {
		fmt.Printf("Loading config from speedtest.net\n")
	}
	c, err := tester.Client.GetConfig()
	if err != nil {
		return err
	}
	tester.Client.Config = &c

	if tester.Debug {
		fmt.Printf("\n")
	}

	if tester.Debug {
		fmt.Printf("Getting servers list...")
	}
	allServers, err := tester.Client.GetServers()
	if err != nil {
		log.Fatal(err)
	}
	if tester.Debug {
		fmt.Printf("(%d) found\n", len(allServers))
	}
	for s := range allServers {
		server := allServers[s]
		print.Server(server)
	}
	return nil
}
