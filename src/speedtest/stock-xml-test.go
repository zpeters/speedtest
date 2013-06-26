package main
import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/xml"
	"time"
	"strconv"
	"sort"
	"math"
	"strings"
	"math/rand"
	_"encoding/hex"
)

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers.php"
var DEBUG = true
var CONFIG Config
const rEarth = 6372.8

// add some debugging timing to different functions or look into profiling


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

func urandom(n int) []byte {
	bytes := make([]byte, n)
	for i:=0; i<n; i++ {
		bytes[i] = byte(rand.Int31())
	}
	return bytes
}


func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if DEBUG { log.Printf("Debugging on...\n") }
	CONFIG = getConfig()

	if DEBUG { log.Printf("Me (%s) - IP: %s - %f,%f\n", CONFIG.Isp, CONFIG.Ip, CONFIG.Lat, CONFIG.Lon) }
	
	allServers := getServers()

	closestServers := getClosestServers(10, allServers)
	if DEBUG {
		for s := range closestServers {
			log.Printf("%s (%s) - %f %f - %f km\n", closestServers[s].Country, closestServers[s].Name , closestServers[s].Lat, closestServers[s].Lon, closestServers[s].Distance)
		}
	}

	fastestServer := getFastestServer(10, closestServers)
	fmt.Printf("Fastest Server: %v\n", fastestServer)


	//Test download
	//Create URLS
	dlsizes := []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
	var urls []string
	for size := range dlsizes {
		for i := 0; i<4; i++ {
			url := fastestServer.Url
			splits := strings.Split(url, "/")
			baseUrl := strings.Join(splits[1:len(splits) -1], "/")
			randomImage := fmt.Sprintf("random%dx%d.jpg", dlsizes[size], dlsizes[size])
			downloadUrl := "http:/" + baseUrl + "/" + randomImage
			
			urls = append(urls, downloadUrl)
		}
	}	

	var speedAcc float64
	for u := range urls {
		fmt.Printf("%s\n", urls[u])
		dlSpeed := downloadSpeed(urls[u])
		fmt.Printf("\tDownload speed: %f Mbps\n", dlSpeed)
		speedAcc = speedAcc + dlSpeed
	}
	fmt.Printf("Average: %f Mbps\n", (speedAcc / float64(len(urls))))

	// // Test upload
	// // https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli
	// ulsizesizes := []int{int(0.25 * 1024 * 1024), int(0.5 * 1024 * 1024)}
	// var ulsize []int
	
	// for size := range ulsizesizes {
	// 	for i := 0; i<25; i++ {
	// 		ulsize = append(ulsize, ulsizesizes[size])
	// 	}
	// }
	// fmt.Printf("Ulsize: %v\n", ulsize)
	// fmt.Printf("Urandom: %v\n", urandom(ulsize[1]))
	// fmt.Printf("Dump: %v\n", hex.Dump(urandom(ulsize[1])))
	
}
