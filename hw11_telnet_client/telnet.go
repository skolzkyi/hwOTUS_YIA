package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
)

var ErrTimeOut = errors.New("timeout error")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	timeout           time.Duration
	inData            io.ReadCloser
	outData           io.Writer
	serviceMessageOut io.Writer
	conn              net.Conn
	addr              string
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tClientExempl := TelnetClientImpl{}
	tClientExempl.timeout = timeout
	tClientExempl.addr = address
	tClientExempl.serviceMessageOut = os.Stderr
	tClientExempl.inData = in
	tClientExempl.outData = out
	message := "New TelnetClient created(addr: " + address + ", timeout: " + timeout.String() + ")\n"
	tClientExempl.serviceMessageOut.Write([]byte(message))

	return &tClientExempl
}

func (tClient *TelnetClientImpl) Connect() error {
	var err error
	tClient.conn, err = net.DialTimeout("tcp", tClient.addr, tClient.timeout)
	if err != nil {
		err = ErrTimeOut
		tClient.serviceMessageOut.Write([]byte(err.Error() + "\n"))
		return err
	}
	tClient.serviceMessageOut.Write([]byte("client connect on addr " + tClient.addr + "\n"))
	return nil
}

func (tClient *TelnetClientImpl) Send() (err error) {
	_, err = io.Copy(tClient.conn, tClient.inData)
	if err != nil {
		tClient.serviceMessageOut.Write([]byte(err.Error() + "\n"))
		return err
	}
	tClient.serviceMessageOut.Write([]byte("Finished Send\n"))

	return nil
}

func (tClient *TelnetClientImpl) Receive() (err error) {
	_, err = io.Copy(tClient.outData, tClient.conn)
	if err != nil {
		tClient.serviceMessageOut.Write([]byte(err.Error() + "\n"))
		return err
	}
	tClient.serviceMessageOut.Write([]byte("Finished Receive\n"))

	return nil
}

func (tClient *TelnetClientImpl) Close() error {
	if err := tClient.conn.Close(); err != nil {
		tClient.serviceMessageOut.Write([]byte(err.Error() + "\n"))
		return err
	}
	tClient.serviceMessageOut.Write([]byte("Connection closed\n"))

	return nil
}
