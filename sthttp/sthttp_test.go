package sthttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckHTTPSuccess(t *testing.T) {
	resp := http.Response{}
	resp.StatusCode = 200
	r := checkHTTP(&resp)
	if r != true {
		t.Fail()
	}
}

func TestCheckHTTPFail(t *testing.T) {
	resp := http.Response{}
	resp.StatusCode = 404
	r := checkHTTP(&resp)
	if r != false {
		t.Fail()
	}
}

func TestGetLatencyURL(t *testing.T) {
	s := Server{}
	s.URL = "http://example.com/speedtest/"
	u := GetLatencyURL(s)
	if u != "http://example.com/speedtest/latency.txt" {
		t.Logf("Got latency URL: %s\n", u)
		t.Fail()
	}
}

func TestServerDistance(t *testing.T) {
	s1 := Server{}
	s1.Distance = 10
	s2 := Server{}
	s2.Distance = 20
	s3 := Server{}
	s3.Distance = 200
	s4 := Server{}
	s4.Distance = 100

	servers := []Server{s3, s4, s2, s1}
	sort.Sort(ByDistance(servers))

	assert.EqualValues(t, servers[0].Distance, 10, "Servers list not sorted by distance")
	assert.EqualValues(t, servers[1].Distance, 20, "Servers list not sorted by distance")
	assert.EqualValues(t, servers[2].Distance, 100, "Servers list not sorted by distance")
	assert.EqualValues(t, servers[3].Distance, 200, "Servers list not sorted by distance")
}

func TestServerLatency(t *testing.T) {
	s1 := Server{}
	s1.Latency = 10
	s2 := Server{}
	s2.Latency = 20
	s3 := Server{}
	s3.Latency = 200
	s4 := Server{}
	s4.Latency = 100

	servers := []Server{s3, s4, s2, s1}
	sort.Sort(ByLatency(servers))

	assert.EqualValues(t, servers[0].Latency, 10, "Servers list not sorted by latency")
	assert.EqualValues(t, servers[1].Latency, 20, "Servers list not sorted by latency")
	assert.EqualValues(t, servers[2].Latency, 100, "Servers list not sorted by latency")
	assert.EqualValues(t, servers[3].Latency, 200, "Servers list not sorted by latency")
}

func TestGetConfig(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_config.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_config.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	c, err := GetConfig(ts.URL)
	if err != nil {
		t.Logf("Cannot get config")
		t.Fatal(err)
	}

	assert.EqualValues(t, c.IP, "23.124.0.25", "IP Doesn't match")
	assert.EqualValues(t, c.Lat, 32.5155, "Latitude doesn't match")
	assert.EqualValues(t, c.Lon, -90.1118, "Longitude doesn't match")
	assert.EqualValues(t, c.Isp, "AT&T U-verse", "ISP Doesn't match")
}

func TestGetConfigNoConnection(t *testing.T) {
	_, err := GetConfig("fail")
	assert.Error(t, err, "An error was expected")
}

func TestGetServers(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	servers, err := GetServers(ts.URL, "")
	if err != nil {
		t.Logf("Cannot get servers")
		t.Fatal(err)
	}

	//sthttp_test.go:127: Server 0: sthttp.Server{URL:"http://88.84.191.230/speedtest/upload.php", Lat:70.0733, Lon:29.7497, Name:"Vadso", Country:"Norway", CC:"NO", Sponsor:"Varanger KraftUtvikling AS", ID:"4600", Distance:0, Latency:0}
	expectURL := "http://88.84.191.230/speedtest/upload.php"
	assert.Equal(t, servers[0].URL, expectURL, fmt.Sprintf("Server 0 url should be: '%s'\n", expectURL))

	expectLat := 59.8833
	assert.Equal(t, servers[100].Lat, expectLat, fmt.Sprintf("Server 10 lat should be: '%s'\n", expectLat))

	expectLon := 15.2
	assert.Equal(t, servers[1005].Lon, expectLon, fmt.Sprintf("Server 1050 lat should be: '%s'\n", expectLat))

	expectName := "Chirchiq"
	assert.Equal(t, servers[2021].Name, expectName, fmt.Sprintf("Server 2021 name should be: '%s'\n", expectName))

	expectCountry := "Lao PDR"
	assert.Equal(t, servers[3321].Country, expectCountry, fmt.Sprintf("Server 3321 name should be: '%s'\n", expectCountry))

	expectCC := "US"
	assert.Equal(t, servers[2222].CC, expectCC, fmt.Sprintf("Server 2222 name should be: '%s'\n", expectCC))

	expectSponsor := "SRT Communications"
	assert.Equal(t, servers[1234].Sponsor, expectSponsor, fmt.Sprintf("Server 1234 name should be: '%s'\n", expectSponsor))

	expectID := "2804"
	assert.Equal(t, servers[666].ID, expectID, fmt.Sprintf("Server 666 name should be: '%s'\n", expectID))

	expectDistance := 0
	assert.EqualValues(t, servers[1].Distance, expectDistance, fmt.Sprintf("Server 1 name should be: '%s'\n", expectDistance))

	expectLatency := 0
	assert.EqualValues(t, servers[21].Latency, expectLatency, fmt.Sprintf("Server 21 name should be: '%s'\n", expectLatency))

}

func TestGetServersBlacklist(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	serversBlacklist, err := GetServers(ts.URL, "3484")
	if err != nil {
		t.Logf("Cannot get servers")
		t.Fatal(err)
	}
	serversAll, err := GetServers(ts.URL, "")
	if err != nil {
		t.Logf("Cannot get servers")
		t.Fatal(err)
	}

	assert.Equal(t, len(serversAll)-1, len(serversBlacklist), "All servers should be one less than blacklist list")

}

func TestGetClosestServers(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	servers, err := GetServers(ts.URL, "")
	if err != nil {
		t.Logf("Cannot get servers")
		t.Fatal(err)
	}

	lat := 32.5155
	lon := -90.1118

	sorted := GetClosestServers(servers, lat, lon)

	assert.Equal(t, sorted[0].ID, "2630", "Closest server ID should be 2630")
}

func TestGetLatency(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, "Hello World")
	}))
	defer ts.Close()

	s := Server{}
	latency, err := GetLatency(s, ts.URL, 5)
	assert.NoError(t, err, "Error getting latency")
	assert.True(t, latency > 100, "Latency faster than expected")
}

func TestGetFastestServer(t *testing.T) {
	x, err := ioutil.ReadFile("sthttp_test_servers.xml")
	if err != nil {
		t.Logf("Cannot read sthttp_test_servers.xml")
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(x))
	}))
	defer ts.Close()

	servers, err := GetServers(ts.URL, "")
	if err != nil {
		t.Logf("Cannot get servers")
		t.Fatal(err)
	}

	fs := GetFastestServer(servers)
	assert.NotNil(t, fs, "No fastest server returned")
}

func TestDownloadSpeed(t *testing.T) {
	f, err := os.Open("random750x750.jpg")
	assert.NoError(t, err, "Can't open test file")
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err, "Can't read test file")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()

	res, err := DownloadSpeed(ts.URL)
	assert.NoError(t, err, "There should be no error")
	assert.True(t, res > 0, "Download speed should be faster than zero")
}

func TestDownloadSpeedBadUrl(t *testing.T) {
	res, err := DownloadSpeed("http://0.0.0.0")
	assert.Error(t, err, "This should fail")
	assert.EqualValues(t, res, 0, "Failed download, so speed should be 0")
}

func TestUploadSpeed(t *testing.T) {
	f, err := os.Open("random750x750.jpg")
	assert.NoError(t, err, "Can't open test file")
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err, "Can't read test file")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, b)
	}))
	defer ts.Close()

	res, err := UploadSpeed(ts.URL, "text/xml", b)
	assert.True(t, res > 0, "Upload speed should be greater than 0")
	assert.NoError(t, err, "Upload should not error out")
}
