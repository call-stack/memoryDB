package main

import (
	"os"
	"os/signal"

	"github.com/call-stack/inmemdb/internal"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	srv := internal.NewServer()

	select {
	case <-stop:
		srv.Stop()
	}
}
