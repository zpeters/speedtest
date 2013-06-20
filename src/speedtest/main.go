package main
import (
	"time"
	"fmt"
	"math"
	"net/http"
	"io/ioutil"
	"strings"
	"strconv"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
)

// Shamelessly borrowed from: https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"
var SpeedtestServersUrl = "http://www.speedtest.net/speedtest-servers.php"
var config Config

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
	// we ar4444444444e interested in
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

	
	//fmt.Printf("Fastest Server Nodes: %v\n", fastestServerNodes)

	return fastestServerNodes
}


func getBestServer(servers []xml.Node) {
	for server := range servers {
		acc := 0
		url := servers[server].Attribute("url").String()
		// God this is ugly
		splits := strings.Split(url, "/")
		baseUrl := strings.Join(splits[1:len(splits) -1], "/")
		latencyUrl := "http:/" + baseUrl + "/latency.txt"
		fmt.Printf("Latency url: %s\n", latencyUrl)

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
			acc += int(finish.Sub(start))
		} else {
			acc += 3600
		}
		
		fmt.Printf("Acc: %d\n", acc)
	}
	
}


func main() {
	config = getConfig()
	servers := getClosestServers()
	fmt.Printf("5 Closest servers: %v\n", servers)
	//server := getBestServer(servers)
	getBestServer(servers)
	//fmt.Printf("Best server: %v\n", server)


}