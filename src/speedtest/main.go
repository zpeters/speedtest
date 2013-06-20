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

func getClosestServers() []xml.Node {
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
		theirlat, _ := strconv.ParseFloat(serverNodes[node].Attribute("lat").Value(), 64)
		theirlon, _ := strconv.ParseFloat(serverNodes[node].Attribute("lon").Value(), 64)
		mylat, _ := strconv.ParseFloat(config.client.Attribute("lat").Value(), 64)
		mylon, _ := strconv.ParseFloat(config.client.Attribute("lon").Value(), 64)
		
		myloc := Coordinate{lat:mylat, lon:mylon}
		theirloc := Coordinate{lat:theirlat, lon:theirlon}
		distance := getDistance(myloc, theirloc)
		
		serverNodes[node].SetAttr("distance", fmt.Sprintf("%f", distance))
	
	}

	// fake this for now, we'll sort the list by distance later
	fastestServerNodes := serverNodes[:5]

	return fastestServerNodes
}


func getBestServer(servers []xml.Node) xml.Node {
	// something is very wrong with our latency calculation
	
	for server := range servers {
		acc := float64(0)
		url := servers[server].Attribute("url").String()
		// God this is ugly
		splits := strings.Split(url, "/")
		baseUrl := strings.Join(splits[1:len(splits) -1], "/")
		latencyUrl := "http:/" + baseUrl + "/latency.txt"
		if DEBUG { fmt.Printf("Latency url: %s\n", latencyUrl) }
		
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
				acc += float64(finish.Sub(start))
			} else {
				acc += 3600
			}
			
			if DEBUG { fmt.Printf("\tAcc (%d): %f\n", i, acc) }
		}
		avg := acc / 3
		if DEBUG { fmt.Printf("\tAverage: %f\n", avg) }
		servers[server].SetAttr("avglatency", fmt.Sprintf("%f", avg))		
	}
	// we will put in code to find the quickest one eventually
	// for now just grab one
	return servers[0]
	
}

func downloadSpeed(urls []string) time.Duration {
	t0 := time.Now()
	for url := range urls {
		fmt.Printf("Downloading %s...\n", urls[url])
		emptyGet(urls[url])
	}
	t1 := time.Now()
	fmt.Printf("Took: %v\n", t1.Sub(t0))
	return t1.Sub(t0)
}


func emptyGet(url string) {
	resp, err := http.Get(url)
	handleErr(err)
	defer resp.Body.Close()
	
	_, err2 := ioutil.ReadAll(resp.Body)
	handleErr(err2)
}



func main() {
	// seed our rng
	//rand.Seed( time.Now().UTC().UnixNano())

	config = getConfig()
	servers := getClosestServers()
	bestServer := getBestServer(servers)
	fmt.Printf("Testing from: %s, %s - %s km away - %s ms latency\n", bestServer.Attribute("sponsor"), bestServer.Attribute("country"), bestServer.Attribute("distance"), bestServer.Attribute("avglatency"))

	sizes := []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 400}
	fmt.Printf("Sizes: %v\n", sizes)
	
	var urls []string
	
	for size := range sizes {
		for i := 0; i<4; i++ {
			url := bestServer.Attribute("url").String()
			splits := strings.Split(url, "/")
			baseUrl := strings.Join(splits[1:len(splits) -1], "/")
			randomImage := fmt.Sprintf("random%dx%d.jpg", sizes[size], sizes[size])
			downloadUrl := "http:/" + baseUrl + "/" + randomImage

			urls = append(urls, downloadUrl)
		}
	}
	fmt.Printf("Urls: %v\n", urls)
	speed := downloadSpeed(urls)
	fmt.Printf("Took: %v\n", speed)
	
}