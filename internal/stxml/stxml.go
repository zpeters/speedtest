package stxml

import (
	"encoding/xml"
)

// TheClient is our users information
type TheClient struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

// XMLConfigSettings is a container for settings
type XMLConfigSettings struct {
	XMLName xml.Name  `xml:"settings"`
	Client  TheClient `xml:"client"`
}

// XMLServer is a candidate server
type XMLServer struct {
	XMLName xml.Name `xml:"server"`
	URL     string   `xml:"url,attr"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
	Name    string   `xml:"name,attr"`
	Country string   `xml:"country,attr"`
	CC      string   `xml:"cc,attr"`
	Sponsor string   `xml:"sponsor,attr"`
	ID      string   `xml:"id,attr"`
}

// TheServersContainer is a list of servers
type TheServersContainer struct {
	XMLName    xml.Name    `xml:"servers"`
	XMLServers []XMLServer `xml:"server"`
}

// ServerSettings is the servers part of the setings
type ServerSettings struct {
	XMLName          xml.Name            `xml:"settings"`
	ServersContainer TheServersContainer `xml:"servers"`
}
