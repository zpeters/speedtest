package server

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Server struct {
	Url string `json:"url"`
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	Distance int `json:"distance"`
	Name string `json:"name"`
	Country string `json:"country"`
	Cc string `json:"cc"`
	Sponsor string `json:"sponsor"`
	Id string `json:"id"`
	Preferred int `json:"preferred"`
	Host string `json:"host"`
	ForcePingSelect int `json:"force_ping_select"`
}

func GetAllServers() (servers []Server) {
	res, err := http.Get("http://www.speedtest.net/api/js/servers?engine=js")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(body, &servers)
	if(err != nil){
		panic(err.Error())
	}

	return servers
}

func GetBestServer() (s Server) {
	// TODO right now we are just picking the first server
	// eventually we need a better algorithm
	servers := GetAllServers()
	return servers[0]
}
