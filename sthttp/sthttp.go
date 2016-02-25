package sthttp

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

import (
	"github.com/dchest/uniuri"
)

import (
	"github.com/zpeters/speedtest/coords"
	"github.com/zpeters/speedtest/debug"
	"github.com/zpeters/speedtest/misc"
	"github.com/zpeters/speedtest/settings"
	"github.com/zpeters/speedtest/stxml"
)

// SpeedtestConfigURL is where we pull the global 'config' from speedtest.net
// Per #39 in some situations we need to have a ?=RANDOMNUMBER or the resulting page is blank
var SpeedtestConfigURL = "http://c.speedtest.net/speedtest-config.php?x=" + uniuri.New()

// SpeedtestServersURL is the global list of speedtest servers
// Per #39 in some situations we need to have a ?=RANDOMNUMBER or the resulting page is blank
var SpeedtestServersURL = "http://c.speedtest.net/speedtest-servers-static.php?x=" + uniuri.New()

// HTTPConfigTimeout is how long we'll wait for a config download to timeout
var HTTPConfigTimeout = time.Duration(15 * time.Second)

// HTTPLatencyTimeout is how long we'll wait for a ping to timeout
var HTTPLatencyTimeout = time.Duration(15 * time.Second)

// HTTPDownloadTimeout is how long we'll wait for a download to timeout
var HTTPDownloadTimeout = time.Duration(15 * time.Minute)

// CONFIG is our global config space
var CONFIG Config

// Config struct holds our config (users current ip, lat, lon and isp)
type Config struct {
	IP  string
	Lat float64
	Lon float64
	Isp string
}

// Server struct is a speedtest candidate server
type Server struct {
	URL      string
	Lat      float64
	Lon      float64
	Name     string
	Country  string
	CC       string
	Sponsor  string
	ID       string
	Distance float64
	Latency  float64
}

// ByDistance allows us to sort servers by distance
type ByDistance []Server

func (server ByDistance) Len() int {
	return len(server)
}

func (server ByDistance) Less(i, j int) bool {
	return server[i].Distance < server[j].Distance
}

func (server ByDistance) Swap(i, j int) {
	server[i], server[j] = server[j], server[i]
}

// ByLatency allows us to sort servers by latency
type ByLatency []Server

func (server ByLatency) Len() int {
	return len(server)
}

func (server ByLatency) Less(i, j int) bool {
	return server[i].Latency < server[j].Latency
}

func (server ByLatency) Swap(i, j int) {
	server[i], server[j] = server[j], server[i]
}

// checkHTTP tests if http response is successful (200) or not
func checkHTTP(resp *http.Response) bool {
	var ok bool
	if resp.StatusCode != 200 {
		ok = false
	} else {
		ok = true
	}
	return ok
}

// GetConfig downloads the master config from speedtest.net
func GetConfig() Config {
	client := &http.Client{
		Timeout: HTTPConfigTimeout,
	}
	req, _ := http.NewRequest("GET", SpeedtestConfigURL, nil)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: 'Could not create connection'\n")
	}
	defer resp.Body.Close()
	if checkHTTP(resp) != true {
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
	c.IP = cx.Client.IP
	c.Lat = misc.ToFloat(cx.Client.Lat)
	c.Lon = misc.ToFloat(cx.Client.Lon)
	c.Isp = cx.Client.Isp

	return *c
}

// GetServers will get the full server list
func GetServers() []Server {
	var servers []Server

	client := &http.Client{
		Timeout: HTTPConfigTimeout,
	}
	req, _ := http.NewRequest("GET", SpeedtestServersURL, nil)
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
		server.URL = s.ServersContainer.XMLServers[xmlServer].URL
		server.Lat = misc.ToFloat(s.ServersContainer.XMLServers[xmlServer].Lat)
		server.Lon = misc.ToFloat(s.ServersContainer.XMLServers[xmlServer].Lon)
		server.Name = s.ServersContainer.XMLServers[xmlServer].Name
		server.Country = s.ServersContainer.XMLServers[xmlServer].Country
		server.CC = s.ServersContainer.XMLServers[xmlServer].CC
		server.Sponsor = s.ServersContainer.XMLServers[xmlServer].Sponsor
		server.ID = s.ServersContainer.XMLServers[xmlServer].ID
		servers = append(servers, *server)
	}
	return servers
}

// GetClosestServers takes the full server list and sorts by distance
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

func getLatencyURL(server Server) string {
	u := server.URL
	splits := strings.Split(u, "/")
	baseURL := strings.Join(splits[1:len(splits)-1], "/")
	latencyURL := "http:/" + baseURL + "/latency.txt"
	return latencyURL
}

// GetLatency will test the latency (ping) the given server NUMLATENCYTESTS times and return either the lowest or average depending on what algorithm is set
func GetLatency(server Server) float64 {
	var latency time.Duration
	var minLatency time.Duration
	var avgLatency time.Duration

	for i := 0; i < settings.NUMLATENCYTESTS; i++ {
		var failed bool
		var finish time.Time

		latencyURL := getLatencyURL(server)
		if debug.DEBUG {
			log.Printf("Testing latency: %s (%s)\n", server.Name, server.Sponsor)
		}

		start := time.Now()

		client := &http.Client{
			Timeout: HTTPLatencyTimeout,
		}
		req, _ := http.NewRequest("GET", latencyURL, nil)
		req.Header.Set("Cache-Control", "no-cache")
		resp, err := client.Do(req)

		if err != nil {
			log.Printf("Cannot test latency of '%s' - 'Cannot contact server'\n", latencyURL)
			failed = true
		} else {
			defer resp.Body.Close()
			finish = time.Now()
			_, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				log.Printf("Cannot test latency of '%s' - 'Cannot read body'\n", latencyURL)
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

		if settings.ALGOTYPE == "max" {
			if minLatency == 0 {
				minLatency = latency
			} else if latency < minLatency {
				minLatency = latency
			}
		} else {
			avgLatency = avgLatency + latency
		}

	}

	if settings.ALGOTYPE == "max" {
		return float64(time.Duration(minLatency.Nanoseconds())*time.Nanosecond) / 1000000
	}
	return float64(time.Duration(avgLatency.Nanoseconds())*time.Nanosecond) / 1000000 / float64(settings.NUMLATENCYTESTS)

}

// GetFastestServer test all servers until we find numServers that
// respond, then find the fastest of them.  Some servers show up in the
// master list but timeout or are "corrupt" therefore we bump their
// latency to something really high (1 minute) and they will drop out of
// this test
func GetFastestServer(servers []Server) Server {
	var successfulServers []Server

	for server := range servers {
		if debug.DEBUG {
			log.Printf("Doing %d runs of %v\n", settings.NUMCLOSEST, servers[server])
		}
		Latency := GetLatency(servers[server])

		if debug.DEBUG {
			log.Printf("Total runs took: %v\n", Latency)
		}

		if Latency > float64(time.Duration(1*time.Minute)) {
			if debug.DEBUG {
				log.Printf("Server %d was too slow, skipping...\n", server)
			}
		} else {
			if debug.DEBUG {
				log.Printf("Server latency was ok %f adding to successful servers list", Latency)
			}
			successfulServers = append(successfulServers, servers[server])
			successfulServers[server].Latency = Latency
		}

		if len(successfulServers) == settings.NUMCLOSEST {
			break
		}
	}

	sort.Sort(ByLatency(successfulServers))
	if debug.DEBUG {
		log.Printf("Server: %v is the fastest server\n", successfulServers[0])
	}
	return successfulServers[0]
}

//Use fix buffer to calculate the length of body
func respBodyLen(resp *http.Response) int {
	l := 0
	buf := make([]byte, 4096)
	for {
		if n, err := resp.Body.Read(buf); err != nil {
			break
		} else {
			l += n
		}
	}

	return l
}

// DownloadSpeed measures the mbps of downloading a URL
func DownloadSpeed(url string) float64 {
	start := time.Now()
	if debug.DEBUG {
		log.Printf("Starting test at: %s\n", start)
	}
	client := &http.Client{
		Timeout: HTTPDownloadTimeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Cannot test download speed of '%s' - 'Cannot contact server'\n", url)
	}
	defer resp.Body.Close()

	bodyLen := respBodyLen(resp)
	finish := time.Now()

	bits := float64(bodyLen * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / float64(seconds)
	return mbps
}

// UploadSpeed measures the mbps to http.Post to a URL
func UploadSpeed(url string, mimetype string, data []byte) float64 {
	buf := bytes.NewBuffer(data)

	start := time.Now()
	if debug.DEBUG {
		log.Printf("Starting test at: %s\n", start)
		log.Printf("Starting test at: %d (nano)\n", start.UnixNano())
	}

	resp, err := http.Post(url, mimetype, buf)
	finish := time.Now()
	if err != nil {
		log.Fatalf("Cannot test upload speed of '%s' - 'Cannot contact server'\n", url)
	}
	defer resp.Body.Close()
	_, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatalf("Cannot test upload speed of '%s' - 'Cannot read body'\n", url)
	}

	if debug.DEBUG {
		log.Printf("Finishing test at: %s\n", finish)
		log.Printf("Finishing test at: %d (nano)\n", finish.UnixNano())
		log.Printf("Took: %d (nano)\n", finish.Sub(start).Nanoseconds())
	}

	bits := float64(len(data) * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / float64(seconds)
	return mbps
}
