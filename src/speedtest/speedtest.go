package main
import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/xml"
	"time"
	"os"
	"strconv"
	"sort"
	"math"
	"strings"
	"math/rand"
	"bytes"
	"flag"
)

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers.php"
var DEBUG = false
var CONFIG Config
var NUMCLOSEST = 5
var NUMLATENCYTESTS = 5
var VERSION = "0.01"
const rEarth = 6372.8

type Coordinate struct {
	lat float64
	lon float64
}

type pos struct {
    φ float64 // latitude, radians
    ψ float64 // longitude, radians
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

type Config struct {
	Ip string
	Lat float64
	Lon float64
	Isp string
}
	
type TheClient struct {
	Ip string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

type XMLConfigSettings struct {
	XMLName xml.Name `xml:"settings"`
	Client TheClient `xml:"client"`
}


type XMLServer struct {
	XMLName xml.Name `xml:"server"`
	Url     string `xml:"url,attr"`
	Lat     string `xml:"lat,attr"`
	Lon     string `xml:"lon,attr"`
	Name    string `xml:"name,attr"`
	Country string `xml:"country,attr"`
	CC      string `xml:"cc,attr"`
	Sponsor string `xml:"sponsor,attr"`
	Id      string `xml:"id,attr"`
}

type TheServersContainer struct {
	XMLName xml.Name `xml:"servers"`
	XMLServers []XMLServer `xml:"server"`
}

type ServerSettings struct {
	XMLName xml.Name `xml:"settings"`
	ServersContainer TheServersContainer `xml:"servers"`
}



// Great Circle
// http://rosettacode.org/wiki/Haversine_formula#Go
func haversine(θ float64) float64 {
    return .5 * (1 - math.Cos(θ))
}

func degPos(lat, lon float64) pos {
    return pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

func hsDist(p1, p2 pos) float64 {
    return 2 * rEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
        math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
}


// Simple error handling
func e(err error) {
	if err != nil {
		log.Panicf("Error: %s\n", err)
	}
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

// shortcut to parse float
func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// Download config from speedtest.net
func getConfig() Config {
	if DEBUG { log.Printf("Getting config...\n") }
	resp, err := http.Get(SpeedtestConfigUrl)
	e(err)
	defer resp.Body.Close()
	if checkHttp(resp) != true {
		log.Panicf("Fail: %s\n", resp.Status)
	}
	
	body, err2 := ioutil.ReadAll(resp.Body)
	e(err2)

	cx := new(XMLConfigSettings)
	
	err3 := xml.Unmarshal(body, &cx)
	e(err3)

	c := new(Config)
	c.Ip = cx.Client.Ip
	c.Lat = toFloat(cx.Client.Lat)
	c.Lon = toFloat(cx.Client.Lon)
	c.Isp = cx.Client.Isp

	return *c
}

// Download server list from speedtest.net
func getServers() []Server {
	var servers []Server

	if DEBUG { log.Printf("Getting servers...\n") }

	resp, err := http.Get(SpeedtestServersUrl)
	e(err)
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	e(err2)

	s := new(ServerSettings)
	
	err3 := xml.Unmarshal(body, &s)
	e(err3)
	
	for xmlServer := range s.ServersContainer.XMLServers {
		server := new(Server)
		server.Url = s.ServersContainer.XMLServers[xmlServer].Url
		server.Lat = toFloat(s.ServersContainer.XMLServers[xmlServer].Lat)
		server.Lon = toFloat(s.ServersContainer.XMLServers[xmlServer].Lon)
		server.Name = s.ServersContainer.XMLServers[xmlServer].Name
		server.Country = s.ServersContainer.XMLServers[xmlServer].Country
		server.CC = s.ServersContainer.XMLServers[xmlServer].CC
		server.Sponsor = s.ServersContainer.XMLServers[xmlServer].Sponsor
		server.Id = s.ServersContainer.XMLServers[xmlServer].Id
		servers = append(servers, *server)
	}

	if DEBUG { log.Printf("Found %d servers...\n", len(servers)) }
	return servers
}


func getClosestServers(numServers int, servers []Server) []Server {
	if DEBUG{ log.Printf("Finding %d closest servers...\n", numServers) }
	// calculate all servers distance from us and save them
	for server := range servers {
		theirlat := servers[server].Lat
		theirlon := servers[server].Lon
		mylat    := CONFIG.Lat
		mylon    := CONFIG.Lon

		theirCoords := Coordinate{lat:theirlat, lon:theirlon}
		myCoords := Coordinate{lat:mylat, lon:mylon}

		servers[server].Distance = hsDist(degPos(myCoords.lat, myCoords.lon), degPos(theirCoords.lat, theirCoords.lon))
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

func getFastestServer(numRuns int, servers []Server) Server {
	for server := range servers {
		var latencyAcc time.Duration
		latencyUrl := getLatencyUrl(servers[server])
		if DEBUG { log.Printf("Testing latency: %s (%s)\n", servers[server].Name, servers[server].Sponsor) }

		for i := 0; i < numRuns; i++ {
			start := time.Now()
			resp, err := http.Get(latencyUrl)
			e(err)
			defer resp.Body.Close()
			
			content, err2 := ioutil.ReadAll(resp.Body)
			e(err2)
			finish := time.Now()
			
			if strings.TrimSpace(string(content)) == "test=test" {
				if DEBUG { fmt.Printf("\tRun %d took: %v\n", i, finish.Sub(start)) }
				latencyAcc = latencyAcc + finish.Sub(start)
			}
		}
		if DEBUG { log.Printf("Total runs took: %v\n", latencyAcc) }
		servers[server].AvgLatency = time.Duration(latencyAcc.Nanoseconds() / int64(numRuns)) * time.Nanosecond
	}

	sort.Sort(ByLatency(servers))
		
	return servers[0]
}

func downloadSpeed(url string) float64 {
	start := time.Now()
	resp, err := http.Get(url)
	e(err)
	defer resp.Body.Close()
	data, err2 := ioutil.ReadAll(resp.Body)
	e(err2)
	finish := time.Now()
 	megabytes := float64(len(data)) / float64(1024) / float64(1024)
	seconds := finish.Sub(start).Seconds()
	mbps := (megabytes * 8) / float64(seconds)

	return mbps
}

func uploadSpeed(url string, mimetype string, data []byte) float64 {
	start := time.Now()
	buf := bytes.NewBuffer(data)
	resp, err := http.Post(url, mimetype, buf)
	e(err)
	defer resp.Body.Close()
	_, err2 := ioutil.ReadAll(resp.Body)
	e(err2)
	finish := time.Now()
	megabytes := float64(len(data)) / float64(1024) / float64(1024)
	seconds := finish.Sub(start).Seconds()
	mbps := (megabytes * 8) / float64(seconds)

	return mbps
}



func urandom(n int) []byte {
	b := make([]byte, n)
	for i:=0; i<n; i++ {
		b[i] = byte(rand.Int31())
	}

	return b
}

func init() {
	flag.BoolVar(&DEBUG, "d", false, "Turn on debugging")
	verFlag := flag.Bool("v", false, "Display version")
	flag.Parse()
	if *verFlag == true {
		fmt.Printf("%s - Version: %s\n", os.Args[0], VERSION)
		os.Exit(0)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	if DEBUG { log.Printf("Debugging on...\n") }
}

func downloadTest(server Server) float64 {
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
		if DEBUG { fmt.Printf("Download test %d\n", u) }
		dlSpeed := downloadSpeed(urls[u])
		if DEBUG { fmt.Printf("\tDownload speed: %f Mbps\n", dlSpeed) }
		speedAcc = speedAcc + dlSpeed
	}
	

	mbps := (speedAcc / float64(len(urls)))
	return mbps
}


func uploadTest(server Server) float64 {
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
		if DEBUG { fmt.Printf("Ulsize: %d\n", ulsize[i]) }
		r := urandom(ulsize[i])
		ulSpeed := uploadSpeed(server.Url, "text/xml", r)
		if DEBUG { fmt.Printf("\tUpload speed: %f Mbps\n", ulSpeed) }
		ulSpeedAcc = ulSpeedAcc + ulSpeed
	}
	
	mbps := ulSpeedAcc / float64(len(ulsize))
	return mbps
}


func main() {
	fmt.Printf("Loading config...\n")
	CONFIG = getConfig()

	fmt.Printf("Getting servers list...")
	allServers := getServers()
	fmt.Printf("(%d) found\n", len(allServers))
	
	fmt.Printf("Finding %d closest servers...\n", NUMCLOSEST)
	// add an option for num closest?
	closestServers := getClosestServers(NUMCLOSEST, allServers)
	
	fmt.Printf("Finding fastest server - testing latency %d times...\n", NUMLATENCYTESTS)
	// add an option for num rums, test how many are necessary
	fastestServer := getFastestServer(NUMLATENCYTESTS, closestServers)
	fmt.Printf("Fastest Server: %s (%s - %s) - %s ping \n", fastestServer.Sponsor, fastestServer.Name, fastestServer.Country, fastestServer.AvgLatency)
	
	fmt.Printf("Starting download test...\n")
	dmbps := downloadTest(fastestServer)
	fmt.Printf("Average Download Speed: %f Mbps\n", dmbps)

	
	fmt.Printf("Starting Upload test...\n")
	umbps := uploadTest(fastestServer)
	fmt.Printf("Average Upload Speed: %f Mbps\n", umbps)
 	
}
