package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrTimeOut        = errors.New("timeout error")
	ErrNotExistServer = errors.New("server not exist")
	addr, port        string
	timeout           time.Duration
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1", "net adress for connection")
	flag.StringVar(&port, "port", "", "port for connection")
	flag.DurationVar(&timeout, "limit", 10, "timeout for connection")
}

func main() {
	flag.Parse()

	fulladress := addr + ":" + port
	tClient := NewTelnetClient(fulladress, timeout, os.Stdin, os.Stdout)
	defer tClient.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-quit

		tClient.Close()
		os.Stderr.WriteString("client stopped on: " + addr)
	}()
	tClient.Connect()
	tClient.Receive()
	tClient.Send()
	wg.Wait()
}

// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
