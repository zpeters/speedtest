package print

import (
	"testing"

	"github.com/urfave/cli"
	"github.com/zpeters/speedtest/sthttp"
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
	stc := sthttp.Client{}
	s := sthttp.Server{}
	s.ID = "123"
	s.Sponsor = "Sponsor"
	s.Name = "Name"
	s.Country = "Country"
	ServerReport(&stc, s)
}

func TestEnvironmentReport(t *testing.T) {
	stc := sthttp.Client{
		Config:          &sthttp.Config{},
		SpeedtestConfig: &sthttp.SpeedtestConfig{},
	}
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		EnvironmentReport(&stc)
		return nil
	}
	app.Run([]string{"testing"})
}
