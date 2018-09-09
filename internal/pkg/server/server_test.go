package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllServers(t *testing.T) {
	servers, err := GetAllServers()
	assert.NoError(t, err)
	aServer := servers[0]
	assert.NotEmpty(t, aServer.ID)
	assert.NotEmpty(t, aServer.Host)
}
