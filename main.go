package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/call-stack/inmemdb/internal"
)

func handleConnection(conn net.Conn, db *internal.Database) {
	//do something here
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		fmt.Println(scanner.Text())
	}

}

func main() {
	stop := make(chan os.Signal, 1)

	srv := internal.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}
