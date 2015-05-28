package print

import (
	"fmt"
	"time"
	"log"
	"runtime"
)

import (
	"github.com/codegangsta/cli"
)

import (
	"github.com/zpeters/speedtest/debug"
	"github.com/zpeters/speedtest/sthttp"
)

func PrintServer(server sthttp.Server) {
	fmt.Printf("%-4s | %s (%s, %s)\n", server.Id, server.Sponsor, server.Name, server.Country)
}

func PrintServerReport(server sthttp.Server, reportchar string) {
	fmt.Printf("%s%s%s%s%s(%s,%s)%s", time.Now(), reportchar, server.Id, reportchar, server.Sponsor, server.Name, server.Country, reportchar)
}

func EnvironmentReport(c *cli.Context, numclosest int, numlatencytests int, reportchar string, algotype string) {
	log.Printf("Env Report")
	log.Printf("-------------------------------\n")
	log.Printf("[User Environment]\n")
	log.Printf("Arch: %v\n", runtime.GOARCH)
	log.Printf("OS: %v\n", runtime.GOOS)
	log.Printf("IP: %v\n", sthttp.CONFIG.Ip)
	log.Printf("Lat: %v\n", sthttp.CONFIG.Lat)
	log.Printf("Lon: %v\n", sthttp.CONFIG.Lon)
	log.Printf("ISP: %v\n", sthttp.CONFIG.Isp)
	log.Printf("-------------------------------\n")
	log.Printf("[Settings]\n")
	if c.Bool("debug") {
		log.Printf("Debug (user): %v\n", debug.DEBUG)
	} else {
		log.Printf("Debug (default): %v\n", debug.DEBUG)
	}
	if c.Bool("quiet") {
		log.Printf("Quiet (user): %v\n", debug.QUIET)
	} else {
		log.Printf("Quiet (default): %v\n", debug.QUIET)
	}
	if c.Int("numclosest") == 0 {
		log.Printf("NUMCLOSEST (default): %v\n", numclosest)
	} else {
		log.Printf("NUMCLOSEST (user): %v\n", numclosest)

	}
	if c.Int("numlatency") == 0 {
		log.Printf("NUMLATENCYTESTS (default): %v\n", numlatencytests)
	} else {
		log.Printf("NUMLATENCYTESTS (user): %v\n", numlatencytests)
	}
	if c.String("server") == "" {
		log.Printf("server (default none specified)\n")
	} else {
		log.Printf("server (user): %s\n", c.String("server"))
	}
	if c.String("reportchar") == "" {
		log.Printf("reportchar (default): %s\n", reportchar)
	} else {
		log.Printf("reportchar (user): %s\n", c.String("reportchar"))
	}
	if c.String("algo") == "" {
		log.Printf("algo (default): %s\n", algotype)
	} else {
		log.Printf("algo (user): %s\n", c.String("algo"))
	}
	log.Printf("--------------------------------\n")
	log.Printf("[Mode]\n")
	log.Printf("Report: %v\n", c.Bool("report"))
	log.Printf("List: %v\n", c.Bool("list"))
	log.Printf("Ping: %v\n", c.Bool("Ping"))
	log.Printf("-------------------------------\n")
}

