package main

import (
	"context"
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

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

	fullAdress := addr + ":" + port
	err := runTelnetClient(fullAdress, timeout, os.Stdin, os.Stdout)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

func runTelnetClient(fullAddress string, timeout time.Duration, in io.ReadCloser, out io.Writer) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tClient := NewTelnetClient(fullAddress, timeout, in, out)

	err := tClient.Connect()
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			defer cancel()
			err := tClient.Send()
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
				return
			}
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			defer cancel()
			err := tClient.Receive()
			if err != nil {
				os.Stderr.WriteString(err.Error() + "\n")
				return
			}
		}
	}()

	<-ctx.Done()

	err = tClient.Close()
	if err != nil {
		return err
	}

	return nil
}
