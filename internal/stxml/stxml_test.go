package stxml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMLStructs(t *testing.T) {
	tc := TheClient{}
	assert.NotNil(t, tc, "Can't create TheClient struct")

	xcs := XMLConfigSettings{}
	assert.NotNil(t, xcs, "Can't create XMLConfigSettings struct")

	xs := XMLServer{}
	assert.NotNil(t, xs, "Can't create XMLServer struct")

	ts := TheServersContainer{}
	assert.NotNil(t, ts, "Can't create TheServersContainer struct")

	ss := ServerSettings{}
	assert.NotNil(t, ss, "Can't create ServerSettings struct")
}
