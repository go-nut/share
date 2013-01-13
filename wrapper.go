package main

import (
	"errors"
	"net"
	"time"
)

type lwrap struct {
	li       net.Listener
	connChan chan bool
	wait     chan bool
}

func (l *lwrap) Accept() (wrapper net.Conn, err error) {
	nconn, err := l.li.Accept()
	if nconn != nil {
		wrapper = &wconn{
			conn:     nconn,
			connChan: l.connChan,
		}
		l.connChan <- true
	}
	return
}

func (l *lwrap) Close() error {
	l.wait <- true
	return l.li.Close()
}

func (l *lwrap) Addr() net.Addr {
	return l.li.Addr()
}

func wrapListener(host string, wait, connChan chan bool) (wrapper net.Listener, err error) {
	l, err := net.Listen("tcp", host)
	wrapper = &lwrap{
		li:       l,
		connChan: connChan,
		wait:     wait,
	}
	return
}

// net.Conn wrapper
type wconn struct {
	conn     net.Conn
	connChan chan bool
}

func (c *wconn) Close() error {
	c.connChan <- false
	return c.conn.Close()
}

func (c *wconn) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *wconn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *wconn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *wconn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *wconn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *wconn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *wconn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// Returns an ipv4 addr
func getAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, interf := range interfaces {
		addrs, err := interf.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			IP, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return "", err
			}
			IP = IP.To4()
			if IP != nil && !IP.IsLoopback() {
				return IP.String(), nil
			}
		}
	}
	return "", errors.New("Could not find interface to bind to")
}
