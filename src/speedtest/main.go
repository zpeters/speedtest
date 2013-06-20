package main
import (
	"fmt"
	"math"
	"net/http"
	"io/ioutil"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
)

// Shamelessly borrowed from: https://github.com/sivel/speedtest-cli/blob/master/speedtest-cli

var SpeedtestConfigUrl = "http://www.speedtest.net/speedtest-config.php"

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



func main() {
	origin := Coordinate{lat:50.00, lon:50.00}
	destination := Coordinate{lat:100.00, lon:100.00}
	fmt.Printf("Distance: %f\n", getDistance(origin, destination))

	config := getConfig()

	fmt.Printf("Client: %v\n", config.client)
 	fmt.Printf("Times: %v\n", config.times)
	fmt.Printf("Download: %v\n", config.download)
	fmt.Printf("Upload: %v\n", config.upload)

	fmt.Printf("Client IP: %v\n", config.client.Attribute("ip"))

}