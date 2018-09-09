package cmds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	assert.NotNil(t, conn)
}

func TestVersion(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	ver := Version(conn)
	assert.Contains(t, ver, ".")
}

func TestPing(t *testing.T) {
	var i int64
	var base float64 = 200
	var delta float64 = 200

	conn := Connect("speedtest.tec.com:8080")
	ms := Ping(conn)
	t.Logf("Base: %#v", base)
	t.Logf("Delta: %#v", delta)
	t.Logf("ms: %#v", ms)
	assert.IsType(t, i, ms)
	assert.InDeltaf(t, ms, base, delta, "Delta %f too large, expected within %fms of %fms", ms, delta, base)
}

func TestCalcMs(t *testing.T) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	finish := time.Now()
	ms := calcMs(start, finish)
	t.Logf("MS: %#v", ms)
	assert.InDelta(t, ms, 100, 10)
}

func TestGenerateBytes(t *testing.T) {
	randBytes1 := 10
	randBytes2 := 555
	randBytes3 := 874321
	bytes1 := generateBytes(randBytes1)
	bytes2 := generateBytes(randBytes2)
	bytes3 := generateBytes(randBytes3)
	assert.Equal(t, randBytes1, len(bytes1))
	assert.Equal(t, randBytes2, len(bytes2))
	assert.Equal(t, randBytes3, len(bytes3))
}

func TestDownload(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	bytes1 := 100
	res1 := Download(conn, bytes1)
	bytes2 := 555555
	res2 := Download(conn, bytes2)
	assert.NotNil(t, res1)
	assert.Equal(t, bytes1, res1.Bytes)
	assert.InDelta(t, res1.DurationMs, 100, 100)
	assert.NotNil(t, res2)
	assert.Equal(t, bytes2, res2.Bytes)
	assert.InDelta(t, res2.DurationMs, 2000, 1000)
}

func TestUpload(t *testing.T) {
	conn := Connect("speedtest.tec.com:8080")
	bytes1 := 100
	res1 := Upload(conn, bytes1)
	bytes2 := 555555
	res2 := Upload(conn, bytes2)
	assert.NotNil(t, res1)
	assert.Equal(t, bytes1, res1.Bytes)
	assert.InDelta(t, res1.DurationMs, 100, 100)
	assert.NotNil(t, res2)
	assert.Equal(t, bytes2, res2.Bytes)
	assert.InDelta(t, res2.DurationMs, 2000, 1000)
}
