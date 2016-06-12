package print

import (
	"testing"

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
