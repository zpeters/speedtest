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

	// fmt.Printf("Version: %s\n", app.Version(conn))

	// ping := app.PingTest(conn, 5)
	// fmt.Printf("Ping results: %d ms\n", ping)

	// dl := app.DownloadTest(conn, 4, 1000000)
	// fmt.Printf("Download results: %s mbps\n", dl)

	ul := app.UploadTest(conn, 2, 5000)
	fmt.Printf("Upload results: %s mbps\n", ul)

	app.Quit(conn)
}
