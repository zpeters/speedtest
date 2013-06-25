package main
import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	_"reflect"
	"encoding/xml"
)

func e(err error) {
	if err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func checkHttp(resp *http.Response) bool {
	var ok bool
	if resp.StatusCode != 200 {
		ok = false
	} else {
		ok = true
	}
	return ok
}

func main() {
	resp, err := http.Get("http://www.speedtest.net/speedtest-config.php")
	e(err)
	defer resp.Body.Close()
	if checkHttp(resp) != true {
		log.Panicf("Fail: %s\n", resp.Status)
	}
	
	body, err2 := ioutil.ReadAll(resp.Body)
	e(err2)
	
	type TheClient struct {
		Ip string `xml:"ip,attr"`
		Lat string `xml:"lat,attr"`
		Lon string `xml:"lon,attr"`
		Isp string `xml:"isp,attr"`
	}

	type Settings struct {
		XMLName xml.Name `xml:"settings"`
		Client TheClient `xml:"client"`
	}

	settings := new(Settings)
	
	err3 := xml.Unmarshal(body, &settings)
	e(err3)
	// fmt.Printf("IP: %v\n", settings.Client[0].Ip)
	// fmt.Printf("Lat: %v\n", settings.Client[0].Lat)
	// fmt.Printf("Lon: %v\n", settings.Client[0].Lon)
	// fmt.Printf("Isp: %v\n", settings.Client[0].Isp)
	fmt.Printf("IP: %v\n", settings.Client.Ip)
	fmt.Printf("Lat: %v\n", settings.Client.Lat)
	fmt.Printf("Lon: %v\n", settings.Client.Lon)
	fmt.Printf("Isp: %v\n", settings.Client.Isp)
}
