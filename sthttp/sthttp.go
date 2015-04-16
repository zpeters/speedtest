package sthttp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

import (
	"github.com/zpeters/speedtest/coords"
	"github.com/zpeters/speedtest/debug"
	"github.com/zpeters/speedtest/misc"
	"github.com/zpeters/speedtest/stxml"
)

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers-static.php"
var HTTP_CONFIG_TIMEOUT = time.Duration(5 * time.Second)
var HTTP_LATENCY_TIMEOUT = time.Duration(5 * time.Second)
var HTTP_DOWNLOAD_TIMEOUT = time.Duration(5 * time.Minute)
var CONFIG Config


type Config struct {
	Ip  string
	Lat float64
	Lon float64
	Isp string
}

type Server struct {
	Url        string
	Lat        float64
	Lon        float64
	Name       string
	Country    string
	CC         string
	Sponsor    string
	Id         string
	Distance   float64
	Latency float64
}

// Sort by Distance
type ByDistance []Server

func (this ByDistance) Len() int {
	return len(this)
}

func (this ByDistance) Less(i, j int) bool {
	return this[i].Distance < this[j].Distance
}

func (this ByDistance) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Sort by latency
type ByLatency []Server

func (this ByLatency) Len() int {
	return len(this)
}

func (this ByLatency) Less(i, j int) bool {
	return this[i].Latency < this[j].Latency
}

func (this ByLatency) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Check http response
func checkHttp(resp *http.Response) bool {
	var ok bool
	if resp.StatusCode != 200 {
		ok = false
	} else {
		ok = true
	}
	return ok
}


// Download config from speedtest.net
func GetConfig() Config {
	client := &http.Client{
		Timeout: HTTP_CONFIG_TIMEOUT,
	}
	req, _ := http.NewRequest("GET", SpeedtestConfigUrl, nil)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: 'Could not create connection'\n")
	}
	defer resp.Body.Close()
	if checkHttp(resp) != true {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: '%s'\n", resp.Status)
	}

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: 'Cannot read body'\n")
	}

	cx := new(stxml.XMLConfigSettings)

	err3 := xml.Unmarshal(body, &cx)
	if err3 != nil {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: 'Cannot unmarshal XML'\n")
	}

	c := new(Config)
	c.Ip = cx.Client.Ip
	c.Lat = misc.ToFloat(cx.Client.Lat)
	c.Lon = misc.ToFloat(cx.Client.Lon)
	c.Isp = cx.Client.Isp

	if debug.DEBUG {
		fmt.Printf("Config: %v\n", c)
	}

	return *c
}

func GetServers() []Server {
	var servers []Server

	client := &http.Client{
		Timeout: HTTP_CONFIG_TIMEOUT,
	}
	req, _ := http.NewRequest("GET", SpeedtestServersUrl, nil)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Cannot get servers list from speedtest.net: 'Cannot contact server'\n")
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Cannot get servers list from speedtest.net: 'Cannot read body'\n")
	}

	s := new(stxml.ServerSettings)

	err3 := xml.Unmarshal(body, &s)
	if err3 != nil {
		log.Fatalf("Cannot get servers list from speedtest.net: 'Cannot unmarshal XML'\n")
	}

	for xmlServer := range s.ServersContainer.XMLServers {
		server := new(Server)
		server.Url = s.ServersContainer.XMLServers[xmlServer].Url
		server.Lat = misc.ToFloat(s.ServersContainer.XMLServers[xmlServer].Lat)
		server.Lon = misc.ToFloat(s.ServersContainer.XMLServers[xmlServer].Lon)
		server.Name = s.ServersContainer.XMLServers[xmlServer].Name
		server.Country = s.ServersContainer.XMLServers[xmlServer].Country
		server.CC = s.ServersContainer.XMLServers[xmlServer].CC
		server.Sponsor = s.ServersContainer.XMLServers[xmlServer].Sponsor
		server.Id = s.ServersContainer.XMLServers[xmlServer].Id
		servers = append(servers, *server)
	}
	return servers
}

func GetClosestServers(servers []Server) []Server {
	if debug.DEBUG {
		log.Printf("Sorting all servers by distance...\n")
	}

	mylat := CONFIG.Lat
	mylon := CONFIG.Lon
	myCoords := coords.Coordinate{Lat: mylat, Lon: mylon}
	for server := range servers {
		theirlat := servers[server].Lat
		theirlon := servers[server].Lon
		theirCoords := coords.Coordinate{Lat: theirlat, Lon: theirlon}

		servers[server].Distance = coords.HsDist(coords.DegPos(myCoords.Lat, myCoords.Lon), coords.DegPos(theirCoords.Lat, theirCoords.Lon))
	}

	sort.Sort(ByDistance(servers))

	return servers
}

func getLatencyUrl(server Server) string {
	u := server.Url
	splits := strings.Split(u, "/")
	baseUrl := strings.Join(splits[1:len(splits)-1], "/")
	latencyUrl := "http:/" + baseUrl + "/latency.txt"
	return latencyUrl
}

func GetLatency(server Server, numRuns int, algotype string) float64 {
	var latency time.Duration
	var minLatency time.Duration
	var avgLatency time.Duration

	for i := 0; i < numRuns; i++ {
		var failed bool
		var finish time.Time
		
		latencyUrl := getLatencyUrl(server)
		if debug.DEBUG {
			log.Printf("Testing latency: %s (%s)\n", server.Name, server.Sponsor)
		}


		start := time.Now()

		client := &http.Client{
			Timeout: HTTP_LATENCY_TIMEOUT,
		}
		req, _ := http.NewRequest("GET", latencyUrl, nil)
		req.Header.Set("Cache-Control", "no-cache")
		resp, err := client.Do(req)

		if err != nil {
		 	log.Printf("Cannot test latency of '%s' - 'Cannot contact server'\n", latencyUrl)
		 	failed = true
		} else {
			defer resp.Body.Close()
			finish = time.Now()
			_, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				log.Printf("Cannot test latency of '%s' - 'Cannot read body'\n", latencyUrl)
				failed = true
			}

		}

		 if failed == true {
			latency = 1 * time.Minute
		 } else {
			 latency = finish.Sub(start)
		 }

		if debug.DEBUG {
			log.Printf("\tRun took: %v\n", latency)
		}

		if algotype == "max" {
			if minLatency == 0 {
				minLatency = latency
			} else if latency < minLatency {
				minLatency = latency
			}
		} else {
			avgLatency = avgLatency + latency
		}
		
	}
	if algotype == "max" {
		return float64(time.Duration(minLatency.Nanoseconds())*time.Nanosecond) / 1000000
	} else {
		return float64(time.Duration(avgLatency.Nanoseconds())*time.Nanosecond) / 1000000 / float64(numRuns)
		
	}
}

func GetFastestServer(numServers int, numRuns int, servers []Server, algotype string) Server {
	// test all servers until we find numServers that respond, then
	// find the fastest of them.  Some servers show up in the master list
	// but timeout or are "corrupt" therefore we bump their latency
	// to something really high (1 minute) and they will drop out of this
	// test
	var successfulServers []Server

	
	for server := range servers {
		if debug.DEBUG {
			log.Printf("Doing %v runs of %s\n", numRuns, servers[server])
		}
		Latency := GetLatency(servers[server], numRuns, algotype)

		if debug.DEBUG {
			log.Printf("Total runs took: %v\n", Latency)
		}

		if (Latency > float64(time.Duration(1 * time.Minute))) {
			if debug.DEBUG {
				log.Printf("Server %s was too slow, skipping...\n", server)
			}
		} else {
			if debug.DEBUG {
				log.Printf("Server latency was ok %v adding to successful servers list", Latency)
			}
			successfulServers = append(successfulServers, servers[server])
			successfulServers[server].Latency = Latency
		}

		if len(successfulServers) == numServers {
			break
		}
	}

	sort.Sort(ByLatency(successfulServers))
	if debug.DEBUG {
		log.Printf("Server: %s is the fastest server\n", successfulServers[0])
	}
	return successfulServers[0]
}

func DownloadSpeed(url string) float64 {
	start := time.Now()
	if debug.DEBUG { log.Printf("Starting test at: %s\n", start) }
	client := &http.Client{
		Timeout: HTTP_DOWNLOAD_TIMEOUT,	
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Cannot test download speed of '%s' - 'Cannot contact server'\n", url)
	}
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Cannot test download speed of '%s' - 'Cannot read body'\n", url)
	}
	finish := time.Now()

	bits := float64(len(data) * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / float64(seconds)
	return mbps
}

func UploadSpeed(url string, mimetype string, data []byte) float64 {
	start := time.Now()
	if debug.DEBUG { log.Printf("Starting test at: %s\n", start) }	

	buf := bytes.NewBuffer(data)
	resp, err := http.Post(url, mimetype, buf)
	if err != nil {
		log.Fatalf("Cannot test upload speed of '%s' - 'Cannot contact server'\n", url)
	}
	defer resp.Body.Close()
	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Cannot test upload speed of '%s' - 'Cannot read body'\n", url)
	}
	finish := time.Now()

	if debug.DEBUG { log.Printf("Finishing test at: %s\n", finish) }

	bits := float64(len(data) * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / float64(seconds)
	return mbps
}
