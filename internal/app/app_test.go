package app

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/zpeters/speedtest/internal/pkg/cmds"
)

func TestGetBestServer(t *testing.T) {
	var expectedWorstSpeed int64 = 200

	is := is.New(t)

	best, err := GetBestServer()

	is.NoErr(err)                                    // should be no errors
	is.True(best.BestTestPing <= expectedWorstSpeed) // best server should be better than 100ms
}

func TestCalcMbps(t *testing.T) {
	is := is.New(t)

	start1 := time.Now()
	end1 := start1.Add(time.Second * 1)
	actual1 := CalcMbps(start1, end1, megabitsToBytes(1))

	start2 := time.Now()
	end2 := start2.Add(time.Second * 8)
	actual2 := CalcMbps(start2, end2, megabitsToBytes(15))

	is.Equal(actual1, float64(1))     // 1 megabit in 1 second should be 1/mbps
	is.Equal(actual2, float64(1.875)) // 1 megabit in 1 second should be 1/mbps
}

func TestPingTest(t *testing.T) {
	var min int64 = 30
	var max int64 = 200

	is := is.New(t)

	conn := cmds.Connect("speedtest.tec.com:8080")
	res := PingTest(conn, 3)

	is.True(inRange(min, max, res)) // ping test should be 100
}

func TestDownloadTest(t *testing.T) {
	var numtests int = 3
	var min1 float64 = 1
	var max1 float64 = 100
	//numbytes := []int{1000000} // 1 megabyte
	numbytes := []int{5000000} // 5 megabyte
	//numbytes := []int{10000000} // 10 megabyte

	is := is.New(t)

	conn := cmds.Connect("speedtest.tec.com:8080")

	res1 := DownloadTest(conn, numbytes, numtests)
	t.Logf("Results: %v", res1)
	is.True(inFloatRange(min1, max1, res1)) // download of 10000 bytes should be between 10 and 100 mbps

}

/* Test Helpers */
func megabitsToBytes(mb int) (bytes int) {
	return mb * 125000
}

func inFloatRange(min float64, max float64, res float64) bool {
	if res >= min && res <= max {
		return true
	}
	return false
}

func inRange(min int64, max int64, res int64) bool {
	if res >= min && res <= max {
		return true
	}
	return false
}
