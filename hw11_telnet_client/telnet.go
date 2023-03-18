package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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
	flag.DurationVar(&timeout, "limit", 0, "timeout for connection")
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	inData  io.ReadCloser
	outData io.Writer
	ctx     context.Context
	conn    net.Conn
	//timeout time.Duration
	cancel            context.CancelFunc
	addr              string
	dialer            *net.Dialer
	serviceMessageOut *os.File
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {

	tClientExempl := TelnetClientImpl{}
	tClientExempl.dialer = &net.Dialer{}
	tClientExempl.ctx, tClientExempl.cancel = context.WithTimeout(context.Background(), timeout)
	tClientExempl.addr = address
	tClientExempl.serviceMessageOut = os.Stderr
	tClientExempl.serviceMessageOut.WriteString("New TelnetClient created(addr: " + address + ", timeout: " + timeout.String() + ")")
	tClientExempl.inData = in
	tClientExempl.outData = out
	return tClientExempl
}

func (tClient TelnetClientImpl) Connect() error {
	var err error
	tClient.conn, err = tClient.dialer.DialContext(tClient.ctx, "tcp", tClient.addr)
	if err != nil {
		return err
	}

	return nil
}

func (tClient TelnetClientImpl) Send() error {

	return nil
}

func (tClient TelnetClientImpl) Receive() error {

	return nil
}

func (tClient TelnetClientImpl) Close() error {
	tClient.conn.Close()
	return nil
}

func readRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text()
			log.Printf("From server: %s", text)
		}
	}
	log.Printf("Finished readRoutine")
}

func writeRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup, stdin chan string) {
	defer wg.Done()
	//scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		case str := <-stdin:
			//if !scanner.Scan() {
			//	break OUTER
			//}
			//str := scanner.Text()
			log.Printf("To server %v\n", str)

			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}

	}
	log.Printf("Finished writeRoutine")
}

func stdinScan() chan string {
	out := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		if scanner.Err() != nil {
			close(out)
		}
	}()
	return out
}

func main() {
	flag.Parse()
	fulladress := addr + ":" + port
	tClient := NewTelnetClient(fulladress, timeout, os.Stdin, os.Stdout)

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
	wg.Wait()
}

// P.S. Author's solution takes no more than 50 lines.
