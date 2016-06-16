package print

import (
	"testing"

	"github.com/urfave/cli"
	"github.com/zpeters/speedtest/internal/sthttp"
)

func TestServer(t *testing.T) {
	s := sthttp.Server{}
	s.ID = "123"
	s.Sponsor = "Sponsor"
	s.Name = "Name"
	s.Country = "Country"

	Server(s)
}

func TestServerReport(t *testing.T) {
	s := sthttp.Server{}
	s.ID = "123"
	s.Sponsor = "Sponsor"
	s.Name = "Name"
	s.Country = "Country"
	ServerReport(s)
}

func TestEnvironmentReportDebugQuietOn(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Turn on debugging",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Quiet mode",
		},
	}
	app.Action = func(c *cli.Context) error {
		EnvironmentReport(c)
		return nil
	}
	args := []string{"testing", "-d", "-q"}
	app.Run(args)
}

func TestEnvironmentReportDebugQuietOff(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Turn on debugging",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Quiet mode",
		},
	}
	app.Action = func(c *cli.Context) error {
		EnvironmentReport(c)
		return nil
	}
	args := []string{"testing"}
	app.Run(args)
}
