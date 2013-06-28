package sthttp

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/xml"
	"time"
	"sort"
	"strings"
	"fmt"
	"bytes"
)

import (
	"speedtest/debug"
	"speedtest/misc"
	"speedtest/stxml"
	"speedtest/coords"
)

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers.php"
var CONFIG Config

type Config struct {
	Ip string
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
	AvgLatency time.Duration
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
	return this[i].AvgLatency < this[j].AvgLatency
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
	if debug.DEBUG { log.Printf("Getting config...\n") }
	resp, err := http.Get(SpeedtestConfigUrl)
	misc.E(err)
	defer resp.Body.Close()
	if checkHttp(resp) != true {
		log.Panicf("Fail: %s\n", resp.Status)
	}
	
	body, err2 := ioutil.ReadAll(resp.Body)
	misc.E(err2)

	cx := new(stxml.XMLConfigSettings)
	
	err3 := xml.Unmarshal(body, &cx)
	misc.E(err3)

	c := new(Config)
	c.Ip = cx.Client.Ip
	c.Lat = misc.ToFloat(cx.Client.Lat)
	c.Lon = misc.ToFloat(cx.Client.Lon)
	c.Isp = cx.Client.Isp

	return *c
}

// Download server list from speedtest.net
func GetServers() []Server {
	var servers []Server

	if debug.DEBUG { log.Printf("Getting servers...\n") }

	resp, err := http.Get(SpeedtestServersUrl)
	misc.E(err)
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	misc.E(err2)

	s := new(stxml.ServerSettings)
	
	err3 := xml.Unmarshal(body, &s)
	misc.E(err3)
	
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

	if debug.DEBUG { log.Printf("Found %d servers...\n", len(servers)) }
	return servers
}


func GetClosestServers(numServers int, servers []Server) []Server {
	if debug.DEBUG{ log.Printf("Finding %d closest servers...\n", numServers) }
	// calculate all servers distance from us and save them
	for server := range servers {
		theirlat := servers[server].Lat
		theirlon := servers[server].Lon
		mylat    := CONFIG.Lat
		mylon    := CONFIG.Lon

		theirCoords := coords.Coordinate{Lat:theirlat, Lon:theirlon}
		myCoords := coords.Coordinate{Lat:mylat, Lon:mylon}

		servers[server].Distance = coords.HsDist(coords.DegPos(myCoords.Lat, myCoords.Lon), coords.DegPos(theirCoords.Lat, theirCoords.Lon))
	}
	
	// sort by distance
	sort.Sort(ByDistance(servers))
	
	// return the top X
	return servers[:numServers]
}

func getLatencyUrl(server Server) string {
	u := server.Url
	splits := strings.Split(u, "/")
	baseUrl := strings.Join(splits[1:len(splits) -1], "/")
	latencyUrl := "http:/" + baseUrl + "/latency.txt"
	return latencyUrl
}

func GetFastestServer(numRuns int, servers []Server) Server {
	for server := range servers {
		var latencyAcc time.Duration
		latencyUrl := getLatencyUrl(servers[server])
		if debug.DEBUG { log.Printf("Testing latency: %s (%s)\n", servers[server].Name, servers[server].Sponsor) }

		for i := 0; i < numRuns; i++ {
			start := time.Now()
			resp, err := http.Get(latencyUrl)
			misc.E(err)
			defer resp.Body.Close()
			
			content, err2 := ioutil.ReadAll(resp.Body)
			misc.E(err2)
			finish := time.Now()
			
			if strings.TrimSpace(string(content)) == "test=test" {
				if debug.DEBUG { fmt.Printf("\tRun %d took: %v\n", i, finish.Sub(start)) }
				latencyAcc = latencyAcc + finish.Sub(start)
			}
		}
		if debug.DEBUG { log.Printf("Total runs took: %v\n", latencyAcc) }
		servers[server].AvgLatency = time.Duration(latencyAcc.Nanoseconds() / int64(numRuns)) * time.Nanosecond
	}

	sort.Sort(ByLatency(servers))
		
	return servers[0]
}


func DownloadSpeed(url string) float64 {
	start := time.Now()
	resp, err := http.Get(url)
	misc.E(err)
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	misc.E(err2)
	finish := time.Now()
 	megabytes := float64(len(data)) / float64(1024) / float64(1024)
	seconds := finish.Sub(start).Seconds()
	mbps := (megabytes * 8) / float64(seconds)

	return mbps
}

func UploadSpeed(url string, mimetype string, data []byte) float64 {
	start := time.Now()
	buf := bytes.NewBuffer(data)
	resp, err := http.Post(url, mimetype, buf)
	misc.E(err)
	defer resp.Body.Close()
	_, err2 := ioutil.ReadAll(resp.Body)
	misc.E(err2)
	finish := time.Now()
	megabytes := float64(len(data)) / float64(1024) / float64(1024)
	seconds := finish.Sub(start).Seconds()
	mbps := (megabytes * 8) / float64(seconds)

	return mbps
}
