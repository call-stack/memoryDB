package main

import (
	"os"

	"github.com/call-stack/inmemdb/internal"
)

func main() {
	stop := make(chan os.Signal, 1)

	srv := internal.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}
