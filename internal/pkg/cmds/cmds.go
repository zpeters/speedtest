package cmds
import (
	"fmt"
	"strings"
	"time"
	"net"
	"math/rand"
)
import (
	"github.com/zpeters/speedtest/internal/pkg/comms"
)

func Connect(server string) (conn net.Conn) {
	return comms.Connect(server)
}

func Version(conn net.Conn) (version string) {
	resp := comms.Command(conn, "HI")
	verLine := strings.Split(resp, " ")
	version = verLine[1]
	return version
}
func Ping(conn net.Conn) (result int64) {
	start := time.Now()

	cmdString := fmt.Sprintf("PING %s", start)

	comms.Command(conn, cmdString)

	finish := time.Now()
	diff := finish.Sub(start)
	return diff.Nanoseconds() / 1000000
}

func Download(conn net.Conn, numbytes int) (result float64) {
	start := time.Now()
	cmdString := fmt.Sprintf("DOWNLOAD %d", numbytes)
	comms.Send(conn, cmdString)
	_ = comms.Recv(conn)
	finish := time.Now()

	mbps := calc_mbps(start, finish, numbytes)
	return mbps
}

func Upload(conn net.Conn, numbytes int) (result float64) {
	rand_bytes := generate_bytes(numbytes)

	bytes_string := fmt.Sprintf("%d", len(rand_bytes))
	len_bytes_string := len(bytes_string)
	final_bytes := len_bytes_string + numbytes + len("UPLOAD_0_\n\n")

	cmdString1 := fmt.Sprintf("UPLOAD %d 0", final_bytes)
	cmdString2 := fmt.Sprintf("%s", rand_bytes)

	start := time.Now()
	comms.Send(conn, cmdString1)
	comms.Send(conn, cmdString2)
	_ = comms.Recv(conn)
	finish := time.Now()

	mbps := calc_mbps(start, finish, numbytes)
	return mbps
}


func calc_mbps(start time.Time, finish time.Time, numbytes int) (mbps float64) {
	diff := finish.Sub(start)
	secs := float64(diff.Nanoseconds()) / float64(1000000000)
	megabits := float64(numbytes) / float64(125000)
	mbps = megabits / secs
	return mbps
}

func generate_bytes(numbytes int) (random []byte) {
	random = make([]byte, numbytes)
	_, err := rand.Read(random)
	if err != nil {
		panic(err)
	}
	return random
}
