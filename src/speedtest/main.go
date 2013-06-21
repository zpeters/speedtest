package main
import (
	"time"
	"fmt"
	"math"
	"net/http"
	"io/ioutil"
	_"math/rand"
	"strings"
	"strconv"
	"sort"
	_"reflect"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
)

// Shamelessly borrowed from: https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers.php"
var config Config
var DEBUG = true

type Coordinate struct {
	lat float64
	lon float64
}

type Config struct {
	client xml.Node
	times xml.Node
	download xml.Node
	upload xml.Node
}

type Server struct {
	url string
	lat float64
	lon float64
	name string
	country string
	cc string
	sponsor string
	id string
	distance float64
	avglatency time.Duration
}


type ByDistance []Server

func (this ByDistance) Len() int {
	return len(this)
}

func (this ByDistance) Less(i, j int) bool {
	return this[i].distance < this[j].distance
}

func (this ByDistance) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}


func getDistance(origin Coordinate, destination Coordinate) float64 {
	// Great Circle calculation
	lat1 := origin.lat
	lon1 := origin.lon
	lat2 := destination.lat
	lon2 := destination.lon
	radius := float64(6371)

	dlat := ((lat2-lat1)*math.Pi)/180
	dlon := ((lon2-lon1)*math.Pi)/180

	a := (math.Sin(dlat/2) * math.Sin(dlat/2) +
		math.Cos((lat1*math.Pi)/180) *
		math.Cos((lat2*math.Pi)/180) *
		math.Sin(dlon/2) *
		math.Sin(dlon/2))

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 -a))
	
	d := radius * c

	return d
}

func getConfig() Config {
	// Download the speedtest.net configuration and return only the data
	// we are interested in
	config := Config{}

	resp, err := http.Get(SpeedtestConfigUrl)
	handleErr(err)
	defer resp.Body.Close()
	
	xml, err2 := ioutil.ReadAll(resp.Body)
	handleErr(err2)

	parsedXml, err3 := gokogiri.ParseXml([]byte(xml))
	handleErr(err3)

	root := parsedXml.Root()

	res, err4 := root.Search("//client")
	handleErr(err4)
	config.client = res[0]

	res, err5 := root.Search("//times")
	handleErr(err5)
	config.times = res[0]

	res, err6 := root.Search("//download")
	handleErr(err6)
	config.download = res[0]


	res, err7 := root.Search("//upload")
	handleErr(err7)
	config.upload = res[0]

	return config
}

func getClosestServers() []Server {
	var servers []Server
	
	// get 5 closets servers
	resp, err := http.Get(SpeedtestServersUrl)
	handleErr(err)
	defer resp.Body.Close()

	xml, err2 := ioutil.ReadAll(resp.Body)
	handleErr(err2)

	parsedXml, err3 := gokogiri.ParseXml([]byte(xml))
	handleErr(err3)

	root := parsedXml.Root()

	serverNodes, err4 := root.Search("//server")
	handleErr(err4)
	for node := range serverNodes {
		s := Server{}
		theirlat, _ := strconv.ParseFloat(serverNodes[node].Attribute("lat").Value(), 64)
		theirlon, _ := strconv.ParseFloat(serverNodes[node].Attribute("lon").Value(), 64)
		mylat, _ := strconv.ParseFloat(config.client.Attribute("lat").Value(), 64)
		mylon, _ := strconv.ParseFloat(config.client.Attribute("lon").Value(), 64)
		
		myloc := Coordinate{lat:mylat, lon:mylon}
		theirloc := Coordinate{lat:theirlat, lon:theirlon}

		s.url = serverNodes[node].Attribute("url").String()
		s.lat = theirlat
		s.lon = theirlon
		s.name = serverNodes[node].Attribute("name").String()
		s.country = serverNodes[node].Attribute("country").String()
		s.cc = serverNodes[node].Attribute("cc").String()
		s.sponsor = serverNodes[node].Attribute("sponsor").String()
		s.id = serverNodes[node].Attribute("id").String()
		s.distance = getDistance(myloc, theirloc)
		
		servers = append(servers, s)
	}


	// sort by distance and return top 5
	sort.Sort(ByDistance(servers))
	return servers[:5]
}


func getBestServer(servers []Server) Server {
	// something is very wrong with our latency calculation
	
	for server := range servers {
		var acc time.Duration
		url := servers[server].url
		// God this is ugly
		splits := strings.Split(url, "/")
		baseUrl := strings.Join(splits[1:len(splits) -1], "/")
		latencyUrl := "http:/" + baseUrl + "/latency.txt"
		if DEBUG { fmt.Printf("\tTesting latency: %s (%s)\n", servers[server].name, servers[server].sponsor) }
		
		for i := 0; i < 3; i++ {
			start := time.Now()
			
			latTest, err := http.Get(latencyUrl)
			if err != nil {
				panic(err)
			}
			defer latTest.Body.Close()
		
			content, err2 := ioutil.ReadAll(latTest.Body)
			if err2 != nil {
				panic(err2)
			}
			
			finish := time.Now()
		
			if strings.TrimSpace(string(content)) == "test=test" {
				acc += finish.Sub(start)
			} else {
				acc += 3600
			}
			
			if DEBUG { fmt.Printf("\t\tAcc (%d): %v\n", i, acc) }
		}
		avg := acc / 4
		if DEBUG { fmt.Printf("\t\tAverage: %s\n", avg) }
		servers[server].avglatency = avg		
	}

	// now sort
	sort.Sort(ByDistance(servers))
	return servers[0]
	
}

func downloadSpeed(urls []string) float64 {
	var datalen int
	t0 := time.Now()
	for url := range urls {
		if DEBUG { fmt.Printf("Downloading %s...\n", urls[url]) }
		datalen = datalen + len(getUrl(urls[url]))
	}
	t1 := time.Now()
	return float64(datalen) / t1.Sub(t0).Seconds()
}


func getUrl(url string) []uint8 {
	resp, err := http.Get(url)
	handleErr(err)
	defer resp.Body.Close()
	
	data, err2 := ioutil.ReadAll(resp.Body)
	handleErr(err2)
	
	return data
}



func main() {
	// seed our rng
	//rand.Seed( time.Now().UTC().UnixNano())

	config = getConfig()
	servers := getClosestServers()
	fmt.Printf("Finding closests servers:\n")
	for s := range servers {
		fmt.Printf("\t%s (sponsored by %s) - %6.2fkm\n", servers[s].name, servers[s].sponsor, servers[s].distance)
	}

	bestServer := getBestServer(servers)
	fmt.Printf("Fastest server: %s (sponsored by %s) - %6.2fkm away - %s ms latency\n", bestServer.sponsor, bestServer.country, bestServer.distance, bestServer.avglatency)

	sizes := []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
 	//sizes := []int{350, 500, 750}
 	if DEBUG { fmt.Printf("Sizes: %v\n", sizes) }
	
 	var urls []string
	
	for size := range sizes {
		for i := 0; i<4; i++ {
			url := bestServer.url
			splits := strings.Split(url, "/")
			baseUrl := strings.Join(splits[1:len(splits) -1], "/")
			randomImage := fmt.Sprintf("random%dx%d.jpg", sizes[size], sizes[size])
			downloadUrl := "http:/" + baseUrl + "/" + randomImage

			urls = append(urls, downloadUrl)
		}
	}
	if DEBUG { fmt.Printf("Urls: %v\n", urls) }
	speed := downloadSpeed(urls)
	fmt.Printf("Download Speed: %f Mbit/s\n", (speed / 1000 / 1000) * 8)
}