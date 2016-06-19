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

func TestEnvironmentReport(t *testing.T) {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		EnvironmentReport(c)
		return nil
	}
	app.Run([]string{"testing"})
}
