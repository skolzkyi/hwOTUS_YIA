package main

import (
	"context"
	"io"
	"net"
	"os"
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
	serviceMessageOut io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tClientExempl := TelnetClientImpl{}
	tClientExempl.dialer = &net.Dialer{}
	tClientExempl.ctx, tClientExempl.cancel = context.WithTimeout(context.Background(), timeout)
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
	tClient.conn, err = tClient.dialer.DialContext(tClient.ctx, "tcp", tClient.addr)
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
