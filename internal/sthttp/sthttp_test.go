package sthttp

import (
	"net/http"
	"testing"
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
