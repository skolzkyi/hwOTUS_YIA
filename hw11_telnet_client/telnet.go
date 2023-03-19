package main

import (
	"bufio"
	"context"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	inData            io.ReadCloser
	outData           io.Writer
	ctx               context.Context
	conn              net.Conn
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
	return &tClientExempl
}

func (tClient *TelnetClientImpl) Connect() error {
	var err error
	tClient.conn, err = tClient.dialer.DialContext(tClient.ctx, "tcp", tClient.addr)
	if err != nil {
		return err
	}

	return nil
}

func (tClient *TelnetClientImpl) Send() error {
	var wg sync.WaitGroup
	wg.Add(1)
	//bufReader := bufio.NewReader(tClient.inData)
	go func(wg *sync.WaitGroup, ctx context.Context) {
		defer wg.Done()
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				//bufReader.WriteTo(tClient.conn)
				//fmt.Printf("tClient.conn: %v\n", tClient.conn)
				io.Copy(tClient.conn, tClient.inData)
			}
		}
	}(&wg, tClient.ctx)
	wg.Wait()
	tClient.serviceMessageOut.WriteString("Finished writeRoutine")
	return nil
}

func (tClient *TelnetClientImpl) Receive() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup, ctx context.Context) {
		defer wg.Done()
		//scanner := bufio.NewScanner(tClient.conn)
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				/*
					if !scanner.Scan() {
						tClient.serviceMessageOut.WriteString("CANNOT SCAN")
						break OUTER
					}
					fmt.Println("Receive: ", scanner.Text())
					tClient.outData.Write(scanner.Bytes())
				*/
				io.Copy(tClient.outData, tClient.conn)
			}
		}

	}(&wg, tClient.ctx)
	wg.Wait()
	tClient.serviceMessageOut.WriteString("Finished Receive")
	return nil
}

func (tClient *TelnetClientImpl) Close() error {
	tClient.conn.Close()
	return nil
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

// P.S. Author's solution takes no more than 50 lines.

//C:\REPO\Go\!OTUS\hwOTUS_YIA\hw11_telnet_client
