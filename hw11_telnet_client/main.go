package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ErrTimeOut = errors.New("timeout error")
	timeout    time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", time.Duration(10)*time.Second, "timeout for connection")
}

func main() {
	flag.Parse()
	addr := flag.Arg(0)
	port := flag.Arg(1)

	if addr == "" {
		os.Stderr.WriteString("address is void\n")
		os.Exit(1)
	}

	if port == "" {
		os.Stderr.WriteString("port is void\n")
		os.Exit(1)
	}

	fulladress := addr + ":" + port
	err := runTelnetClient(fulladress, timeout, os.Stdin, os.Stdout)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func runTelnetClient(fulladress string, timeout time.Duration, in io.ReadCloser, out io.Writer) error {
	tClient := NewTelnetClient(fulladress, timeout, in, out)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := tClient.Connect()
	if err != nil {
		return err
	}

	go func() {
		defer cancel()
		err := tClient.Send()
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return
		}
	}()

	go func() {
		defer cancel()
		err := tClient.Receive()
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			return
		}
	}()

	<-ctx.Done()

	err = tClient.Close()
	if err != nil {
		return err
	}

	return nil
}
