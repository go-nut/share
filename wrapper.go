package main

import (
  "net"
  "time"
)

type lwrap struct {
  li net.Listener
}


func (l *lwrap) Accept() (wrapper net.Conn, err error) {
  nconn, err := l.li.Accept()
  wrapper = &wconn {
    conn: nconn,
  }
}


func (l *lwrap) Close() error {
  return l.li.Close()
}

func (l *lwrap) Addr() net.Addr {
  return l.li.Addr()
}


// net.Conn wrapper
type wconn struct {
  conn net.Conn
}

func (c *wconn) Close() error {
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
