package cmds

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

import (
	"github.com/spf13/viper"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/comms"
)

type Result struct {
	Start      time.Time
	Finish     time.Time
	DurationMs int64
	Bytes      int
}

// Connect will returns a socket connection from the server
func Connect(server string) (conn net.Conn) {
	return comms.Connect(server)
}

// Version retrieves and parses the protocol version
func Version(conn net.Conn) (version string) {
	resp, err := comms.Command(conn, "HI")
	if err != nil {
		log.Fatal(err)
	}
	verLine := strings.Split(resp, " ")
	version = verLine[1]
	return version
}

// Ping issues and times a ping command
func Ping(conn net.Conn) (ms int64) {
	start := time.Now()
	cmdString := fmt.Sprintf("PING %s", start)
	comms.Command(conn, cmdString)
	finish := time.Now()
	ms = calcMs(start, finish)
	return ms
}

// Download performs a timed download, returning the mpbs
func Download(conn net.Conn, numbytes int) (result Result) {
	start := time.Now()
	cmdString := fmt.Sprintf("DOWNLOAD %d", numbytes)
	comms.Send(conn, cmdString)
	_, err := comms.Recv(conn)
	if err != nil {
		log.Fatal(err)
	}
	finish := time.Now()

	result = Result{
		Start:      start,
		Finish:     finish,
		DurationMs: calcMs(start, finish),
		Bytes:      numbytes,
	}

	return result
}

// Upload performs a timed upload of numbytes random bytes, returning the mpbs
func Upload(conn net.Conn, numbytes int) (result Result) {
	randBytes := generateBytes(numbytes)

	bytesString := fmt.Sprintf("%d", len(randBytes))
	lenBytesString := len(bytesString)
	finalBytes := lenBytesString + numbytes + len("UPLOAD_0_\n\n")

	cmdString1 := fmt.Sprintf("UPLOAD %d 0", finalBytes)
	cmdString2 := fmt.Sprintf("%s", randBytes)

	start := time.Now()
	comms.Send(conn, cmdString1)
	comms.Send(conn, cmdString2)
	_, err := comms.Recv(conn)
	if err != nil {
		log.Fatal(err)
	}
	finish := time.Now()

	result = Result{
		Start:      start,
		Finish:     finish,
		DurationMs: calcMs(start, finish),
		Bytes:      numbytes,
	}

	return result
}

func generateBytes(numbytes int) (random []byte) {
	if viper.GetBool("Debug") {
		log.Printf("Generating %d random bytes", numbytes)
	}
	random = make([]byte, numbytes)
	_, err := rand.Read(random)
	if err != nil {
		panic(err)
	}
	return random
}

func calcMs(start time.Time, finish time.Time) (ms int64) {
	diff := finish.Sub(start)
	return diff.Nanoseconds() / 1000000
}
