package main
import (
	"fmt"
)
import (
	"github.com/zpeters/speedtest/internal/app"
)

func main() {
	server := app.GetBestServer()
	fmt.Printf("Found best server: (%s) %s - %s\n", server.Id, server.Name, server.Sponsor)
	conn := app.Connect(server.Host)

	fmt.Printf("Speedtest protocol version: %s\n", app.Version(conn))

	ping := app.PingTest(conn, 20)
	download_bytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000, 2000000, 5000000, 7000000}
	download := app.DownloadTest(conn, download_bytes, 4)
	upload_bytes := []int{5000, 10000, 53725, 71582, 73434, 80026, 121474, 1000000, 2000000, 5000000, 7000000}
	upload := app.UploadTest(conn, upload_bytes, 4)

	fmt.Printf("--| Results |---\n")
	fmt.Printf("Ping results: %d ms\n", ping)
	fmt.Printf("Download results: %f mbps\n", download)
	fmt.Printf("Upload results: %f mbps\n", upload)
}
