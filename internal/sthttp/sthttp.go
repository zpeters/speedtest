package sthttp

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/zpeters/speedtest/internal/coords"
	"github.com/zpeters/speedtest/internal/misc"
	"github.com/zpeters/speedtest/internal/stxml"

	"github.com/spf13/viper"
)

// HTTPConfigTimeout is how long we'll wait for a config download to timeout
var HTTPConfigTimeout = time.Duration(viper.GetDuration("httpconfigtimeout") * time.Second)

// HTTPLatencyTimeout is how long we'll wait for a ping to timeout
var HTTPLatencyTimeout = time.Duration(viper.GetDuration("httplatencytimeout") * time.Second)

// HTTPDownloadTimeout is how long we'll wait for a download to timeout
var HTTPDownloadTimeout = time.Duration(viper.GetDuration("httpdownloadtimeout") * time.Minute)

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

// checkBlacklisted tests if the server is on the specified blacklist
func checkBlacklisted(blacklist string, server string) bool {
	var isBlacklisted bool = false
	bl := strings.Split(blacklist, ",")
	for b := range bl {
		blServer := strings.TrimSpace(bl[b])
		if server == blServer {
			isBlacklisted = true
		}
	}
	return isBlacklisted
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
func GetConfig(url string) (c Config, err error) {
	c = Config{}

	client := &http.Client{
		Timeout: HTTPConfigTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Unofficial CLI")

	resp, err := client.Do(req)
	if err != nil {
		return c, err
	}
	defer resp.Body.Close()
	if checkHTTP(resp) != true {
		log.Fatalf("Couldn't retrieve our config from speedtest.net: '%s'\n", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)

	cx := new(stxml.XMLConfigSettings)

	err = xml.Unmarshal(body, &cx)

	c.IP = cx.Client.IP
	c.Lat = misc.ToFloat(cx.Client.Lat)
	c.Lon = misc.ToFloat(cx.Client.Lon)
	c.Isp = cx.Client.Isp

	return c, err
}

// GetServers will get the full server list
func GetServers(url string, blacklist string) (servers []Server, err error) {
	client := &http.Client{
		Timeout: HTTPConfigTimeout,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Unofficial CLI")

	resp, err := client.Do(req)

	if err != nil {
		return servers, err
	}
	defer resp.Body.Close()

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return servers, err2
	}

	s := new(stxml.ServerSettings)

	err3 := xml.Unmarshal(body, &s)
	if err3 != nil {
		return servers, err3
	}

	for xmlServer := range s.ServersContainer.XMLServers {
		// check if server is blacklisted
		if checkBlacklisted(blacklist, s.ServersContainer.XMLServers[xmlServer].ID) == false {
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
	}
	return servers, nil
}

// GetClosestServers takes the full server list and sorts by distance
func GetClosestServers(servers []Server, lat float64, lon float64) []Server {
	if viper.GetBool("debug") {
		log.Printf("Sorting all servers by distance...\n")
	}

	myCoords := coords.Coordinate{Lat: lat, Lon: lon}
	for server := range servers {
		theirlat := servers[server].Lat
		theirlon := servers[server].Lon
		theirCoords := coords.Coordinate{Lat: theirlat, Lon: theirlon}

		servers[server].Distance = coords.HsDist(coords.DegPos(myCoords.Lat, myCoords.Lon), coords.DegPos(theirCoords.Lat, theirCoords.Lon))
	}

	sort.Sort(ByDistance(servers))

	return servers
}

// GetLatencyURL will return the proper url for the latency
// test file when given a server name
func GetLatencyURL(server Server) string {
	u := server.URL
	splits := strings.Split(u, "/")
	baseURL := strings.Join(splits[1:len(splits)-1], "/")
	latencyURL := "http:/" + baseURL + "/latency.txt"
	return latencyURL
}

// GetLatency will test the latency (ping) the given server NUMLATENCYTESTS times and return either the lowest or average depending on what algorithm is set
func GetLatency(server Server, url string, numtests int) (result float64, err error) {
	var latency time.Duration
	var minLatency time.Duration
	var avgLatency time.Duration

	for i := 0; i < numtests; i++ {
		var failed bool
		var finish time.Time

		if viper.GetBool("debug") {
			log.Printf("Testing latency: %s (%s)\n", server.Name, server.Sponsor)
		}

		start := time.Now()

		client, err := getHttpClient()
		if err != nil {
			return result, err
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", "Unofficial CLI")

		resp, err := client.Do(req)

		if err != nil {
			return result, err
		}

		defer resp.Body.Close()
		finish = time.Now()
		_, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return result, err
		}

		if failed == true {
			latency = 1 * time.Minute
		} else {
			latency = finish.Sub(start)
		}

		if viper.GetBool("debug") {
			log.Printf("\tRun took: %v\n", latency)
		}

		if viper.GetString("algotype") == "max" {
			if minLatency == 0 {
				minLatency = latency
			} else if latency < minLatency {
				minLatency = latency
			}
		} else {
			avgLatency = avgLatency + latency
		}

	}

	if viper.GetString("algotype") == "max" {
		result = float64(time.Duration(minLatency.Nanoseconds())*time.Nanosecond) / 1000000
	} else {
		result = float64(time.Duration(avgLatency.Nanoseconds())*time.Nanosecond) / 1000000 / float64(numtests)
	}

	return result, nil

}

// GetFastestServer test all servers until we find numServers that
// respond, then find the fastest of them.  Some servers show up in the
// master list but timeout or are "corrupt" therefore we bump their
// latency to something really high (1 minute) and they will drop out of
// this test
func GetFastestServer(servers []Server) Server {
	var successfulServers []Server

	for server := range servers {
		if viper.GetBool("debug") {
			log.Printf("Doing %d runs of %v\n", viper.GetInt("numclosest"), servers[server])
		}
		Latency, err := GetLatency(servers[server], GetLatencyURL(servers[server]), viper.GetInt("numlatencytests"))
		if err != nil {
			log.Fatal(err)
		}

		if viper.GetBool("debug") {
			log.Printf("Total runs took: %v\n", Latency)
		}

		if Latency > float64(time.Duration(1*time.Minute)) {
			if viper.GetBool("debug") {
				log.Printf("Server %d was too slow, skipping...\n", server)
			}
		} else {
			if viper.GetBool("debug") {
				log.Printf("Server latency was ok %f adding to successful servers list", Latency)
			}
			successfulServers = append(successfulServers, servers[server])
			successfulServers[server].Latency = Latency
		}

		if len(successfulServers) == viper.GetInt("numclosest") {
			break
		}
	}

	sort.Sort(ByLatency(successfulServers))
	if viper.GetBool("debug") {
		log.Printf("Server: %v is the fastest server\n", successfulServers[0])
	}
	return successfulServers[0]
}

// DownloadSpeed measures the mbps of downloading a URL
func DownloadSpeed(url string) (speed float64, err error) {
	start := time.Now()
	if viper.GetBool("debug") {
		log.Printf("Starting test at: %s\n", start)
	}
	client, err := getHttpClient()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Unofficial CLI")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyLen := len(body)
	finish := time.Now()

	bits := float64(bodyLen * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()
	mbps := megabits / float64(seconds)

	return mbps, err
}

// UploadSpeed measures the mbps to http.Post to a URL
func UploadSpeed(url string, mimetype string, data []byte) (speed float64, err error) {
	buf := bytes.NewBuffer(data)

	start := time.Now()
	if viper.GetBool("debug") {
		log.Printf("Starting test at: %s\n", start)
		log.Printf("Starting test at: %d (nano)\n", start.UnixNano())
	}

	client, err := getHttpClient()
	if err != nil {
		return 0, err
	}
	resp, err := client.Post(url, mimetype, buf)
	finish := time.Now()
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if viper.GetBool("debug") {
		log.Printf("Finishing test at: %s\n", finish)
		log.Printf("Finishing test at: %d (nano)\n", finish.UnixNano())
		log.Printf("Took: %d (nano)\n", finish.Sub(start).Nanoseconds())
	}

	bits := float64(len(data) * 8)
	megabits := bits / float64(1000) / float64(1000)
	seconds := finish.Sub(start).Seconds()

	mbps := megabits / float64(seconds)
	return mbps, nil
}

func getSourceIP() (string, error) {
	interfaceOption := viper.GetString("interface")
	if interfaceOption == "" {
		return "", nil
	}

	// does it look like an IP address?
	if net.ParseIP(interfaceOption) != nil {
		return interfaceOption, nil
	}

	// assume that it is the name of an interface
	iface, err := net.InterfaceByName(interfaceOption)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			// fixme: IPv6 support is missing
			if v.IP.To4() != nil {
				return v.IP.String(), nil
			}
		case *net.IPAddr:
			if v.IP.To4() != nil {
				return v.IP.String(), nil
			}
		}
	}

	return "", errors.New("no address found")
}

func getHttpClient() (*http.Client, error) {
	var dialer net.Dialer

	sourceIP, err := getSourceIP()
	if err != nil {
		return nil, err
	}
	if sourceIP != "" {
		bindAddrIP, err := net.ResolveIPAddr("ip", sourceIP)
		if err != nil {
			return nil, err
		}
		bindAddr := net.TCPAddr{
			IP: bindAddrIP.IP,
		}
		dialer = net.Dialer{
			LocalAddr: &bindAddr,
			Timeout:   HTTPConfigTimeout,
			KeepAlive: HTTPConfigTimeout,
		}
	} else {
		dialer = net.Dialer{
			Timeout:   HTTPConfigTimeout,
			KeepAlive: HTTPConfigTimeout,
		}
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: HTTPConfigTimeout,
	}
	client := &http.Client{
		Timeout:   HTTPConfigTimeout,
		Transport: transport,
	}
	return client, nil
}
