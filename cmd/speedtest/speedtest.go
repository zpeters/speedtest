package main
import (
	"fmt"
)
import (
	"github.com/zpeters/speedtest/internal/app"
)

func main() {
	server := "speedtest1.mtaonline.net:8080"
	conn := app.Connect(server)

	fmt.Printf("Version: %s\n", app.Version(conn))

	ping := app.PingTest(conn, 4)
	fmt.Printf("Ping results: %d ms\n", ping)

	dl := app.DownloadTest(conn, 4, 1000000)
	fmt.Printf("Download results: %#v\n", dl)

	app.Quit(conn)
}
