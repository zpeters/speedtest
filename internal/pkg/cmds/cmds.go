package cmds

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/comms"
)

// Connect will returns a socket connection from the server
func Connect(server string) (conn net.Conn) {
	return comms.Connect(server)
}

// Version retrieves and parses the protocol version
func Version(conn net.Conn) (version string) {
	resp := comms.Command(conn, "HI")
	verLine := strings.Split(resp, " ")
	version = verLine[1]
	return version
}

// Ping issues and times a ping command
func Ping(conn net.Conn) (result int64) {
	start := time.Now()

	cmdString := fmt.Sprintf("PING %s", start)

	comms.Command(conn, cmdString)

	finish := time.Now()
	diff := finish.Sub(start)
	return diff.Nanoseconds() / 1000000
}

// Download performs a timed download, returning the mpbs
func Download(conn net.Conn, numbytes int) (mbps float64) {
	start := time.Now()
	cmdString := fmt.Sprintf("DOWNLOAD %d", numbytes)
	comms.Send(conn, cmdString)
	_ = comms.Recv(conn)
	finish := time.Now()

	mbps = calcMbps(start, finish, numbytes)
	return mbps
}

// Upload performs a timed upload of numbytes random bytes, returning the mpbs
func Upload(conn net.Conn, numbytes int) (result float64) {
	randBytes := generateBytes(numbytes)

	bytesString := fmt.Sprintf("%d", len(randBytes))
	lenBytesString := len(bytesString)
	finalBytes := lenBytesString + numbytes + len("UPLOAD_0_\n\n")

	cmdString1 := fmt.Sprintf("UPLOAD %d 0", finalBytes)
	cmdString2 := fmt.Sprintf("%s", randBytes)

	start := time.Now()
	comms.Send(conn, cmdString1)
	comms.Send(conn, cmdString2)
	_ = comms.Recv(conn)
	finish := time.Now()

	mbps := calcMbps(start, finish, numbytes)
	return mbps
}

func calcMbps(start time.Time, finish time.Time, numbytes int) (mbps float64) {
	diff := finish.Sub(start)
	secs := float64(diff.Nanoseconds()) / float64(1000000000)
	megabits := float64(numbytes) / float64(125000)
	mbps = megabits / secs
	return mbps
}

func generateBytes(numbytes int) (random []byte) {
	random = make([]byte, numbytes)
	_, err := rand.Read(random)
	if err != nil {
		panic(err)
	}
	return random
}
