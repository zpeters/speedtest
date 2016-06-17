package sthttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

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
	u := getLatencyURL(s)
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
