package stxml

import (
	"encoding/xml"
)

type TheClient struct {
	Ip  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

type XMLConfigSettings struct {
	XMLName xml.Name  `xml:"settings"`
	Client  TheClient `xml:"client"`
}

type XMLServer struct {
	XMLName xml.Name `xml:"server"`
	Url     string   `xml:"url,attr"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
	Name    string   `xml:"name,attr"`
	Country string   `xml:"country,attr"`
	CC      string   `xml:"cc,attr"`
	Sponsor string   `xml:"sponsor,attr"`
	Id      string   `xml:"id,attr"`
}

type TheServersContainer struct {
	XMLName    xml.Name    `xml:"servers"`
	XMLServers []XMLServer `xml:"server"`
}

type ServerSettings struct {
	XMLName          xml.Name            `xml:"settings"`
	ServersContainer TheServersContainer `xml:"servers"`
}
