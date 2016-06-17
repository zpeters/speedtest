package print

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/zpeters/speedtest/internal/sthttp"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

// Server prints the results in "human" format
func Server(server sthttp.Server) {
	fmt.Printf("%-4s | %s (%s, %s)\n", server.ID, server.Sponsor, server.Name, server.Country)
}

// ServerReport prints restults in a machine useable format
func ServerReport(server sthttp.Server) {
	fmt.Printf("%s%s%s%s%s(%s,%s)%s", time.Now(), viper.GetString("reportchar"), server.ID, viper.GetString("reportchar"), server.Sponsor, server.Name, server.Country, viper.GetString("reportchar"))
}

// EnvironmentReport is a debugging report helpful for debugging
func EnvironmentReport(c *cli.Context) {
	log.Printf("Env Report")
	log.Printf("-------------------------------\n")
	log.Printf("[User Environment]\n")
	log.Printf("Arch: %v\n", runtime.GOARCH)
	log.Printf("OS: %v\n", runtime.GOOS)
	log.Printf("IP: %v\n", sthttp.CONFIG.IP)
	log.Printf("Lat: %v\n", sthttp.CONFIG.Lat)
	log.Printf("Lon: %v\n", sthttp.CONFIG.Lon)
	log.Printf("ISP: %v\n", sthttp.CONFIG.Isp)
	log.Printf("Config: %s\n", viper.GetString("speedtestconfigurl"))
	log.Printf("Servers: %s\n", sthttp.SpeedtestServersURL)
	log.Printf("-------------------------------\n")
	log.Printf("[args]\n")
	log.Printf("%#v\n", os.Args)
	log.Printf("--------------------------------\n")
	log.Printf("[Mode]\n")
	log.Printf("Report: %v\n", c.Bool("report"))
	log.Printf("List: %v\n", c.Bool("list"))
	log.Printf("Ping: %v\n", c.Bool("Ping"))
	log.Printf("-------------------------------\n")
}
