package cmds

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestConnect(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")

	is.True(conn != nil) // conn should not be nil
}

func TestVersion(t *testing.T) {
	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	ver := Version(conn)

	is.True(isVersion(ver)) // ver should be in the format 1.2
}

func TestPing(t *testing.T) {
	var min int64 = 10
	var max int64 = 300

	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	ms := Ping(conn)

	is.True(inRange(ms, min, max)) // ping should be between 10 and 300 ms
}

func TestGenerateBytes(t *testing.T) {
	is := is.New(t)

	randBytes1 := 10
	randBytes2 := 555
	randBytes3 := 874321

	bytes1 := generateBytes(randBytes1)
	bytes2 := generateBytes(randBytes2)
	bytes3 := generateBytes(randBytes3)

	is.Equal(randBytes1, len(bytes1)) // requesting 10 randBytes should return 10 bytes
	is.Equal(randBytes2, len(bytes2)) // requesting 555 randBytes should return 10 bytes
	is.Equal(randBytes3, len(bytes3)) // requesting 874321 randBytes should return 10 bytes
}

func TestDownload(t *testing.T) {
	var min1 int64 = 100
	var max1 int64 = 200
	var min2 int64 = 1000
	var max2 int64 = 4000

	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	bytes1 := 100
	res1 := Download(conn, bytes1)
	bytes2 := 555555
	res2 := Download(conn, bytes2)

	is.True(inRange(res1.DurationMs, min1, max1)) // 10 byte download should take between 100 to 200 ms
	is.True(inRange(res2.DurationMs, min2, max2)) // 555,555 byte download should take between 1000 and 2000 ms
}

func TestUpload(t *testing.T) {
	var min1 int64 = 100
	var max1 int64 = 200
	var min2 int64 = 1000
	var max2 int64 = 4000

	is := is.New(t)

	conn := Connect("speedtest.tec.com:8080")
	bytes1 := 100
	res1 := Upload(conn, bytes1)
	bytes2 := 555555
	res2 := Upload(conn, bytes2)

	is.True(inRange(res1.DurationMs, min1, max1)) // 10 byte upload should take between 100 to 200 ms
	is.True(inRange(res2.DurationMs, min2, max2)) // 555,555 byte upload should take between 1000 and 2000 ms
}

/* Helpers */
func isVersion(ver string) bool {
	v := strings.Split(ver, ".")
	maj := v[0]
	min := v[1]

	_, err1 := strconv.Atoi(maj)
	_, err2 := strconv.Atoi(min)

	if err1 == nil && err2 == nil && len(v) == 2 {
		return true
	}
	return false
}

func inRange(ms int64, min int64, max int64) bool {
	fmt.Printf("ms: %d", ms)
	fmt.Printf("min: %d", min)
	fmt.Printf("max: %d", max)
	if ms >= min && ms <= max {
		return true
	}
	return false
}
