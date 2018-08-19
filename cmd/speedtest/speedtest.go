package main
import (
	"fmt"
)
import (
	"github.com/zpeters/speedtest/internal/app"
)

func main() {
	//server := "speedtest.turk.net:8080"
	server := "speedtest1.mtaonline.net:8080"
	conn := app.Connect(server)
	//fmt.Println(app.Version(conn))
	res := app.PingTest(conn, 10)
	fmt.Printf("Ping results: %d ms\n", res)
	app.Quit(conn)
}
