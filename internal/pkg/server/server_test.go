package server

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetAllServers(t *testing.T) {
	is := is.New(t)

	servers, err := GetAllServers()
	aServer := servers[0]

	is.NoErr(err)               // we should be able to get servers
	is.True(aServer.ID != "")   // a server ID should not be empty
	is.True(aServer.Host != "") // a server Host should not be empty
}
