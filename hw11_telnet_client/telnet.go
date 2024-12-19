package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	input   io.ReadCloser
	output  io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		input:   in,
		output:  out,
	}
}

func (t *client) Close() error {
	if err := t.connected(); err != nil {
		return err
	}

	return t.conn.Close()
}

func (t *client) Send() error {
	if err := t.connected(); err != nil {
		return err
	}

	_, err := io.Copy(t.conn, t.input)

	return err
}

func (t *client) Receive() error {
	if err := t.connected(); err != nil {
		return err
	}

	_, err := io.Copy(t.output, t.conn)

	return err
}

func (t *client) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)

	return err
}

func (t *client) connected() error {
	if t.conn == nil {
		return fmt.Errorf("tcp connection is nil")
	}

	return nil
}
