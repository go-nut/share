package main

import (
	"net/http"
	"net"
	"sync"
  "errors"
)

type StopServer struct {
	http.Server
	listener net.Listener
	connChan chan bool
}




func (srv *StopServer) ListenAndServe() error {
  if srv.Addr == "" {
    srv.Addr = ":http"
  }
  var err error
  // Wrap the listener so we can know when to close
  l, err := net.Listen("tcp", srv.Addr)
  lw := &lwrap {
    li: l,
  }
  srv.listener = lw

  if err != nil {
    return err
  }
  return srv.Serve(srv.listener)
}

func (srv *StopServer) Serve(l net.Listener) error {
    cur_handler := srv.Handler
    defer func() {
      srv.Handler = cur_handler
	}()
    new_handler := http.NewServeMux()
    new_handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		srv.waiter.Add(1)

		defer srv.waiter.Done()
		cur_handler.ServeHTTP(w, r)
	})
	srv.Handler = new_handler
	return srv.Server.Serve(l)
}

func (srv *StopServer) Stop() error {
	return srv.listener.Close()
}

func (srv *StopServer) WaitUnfinished() {
	srv.waiter.Wait()
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
      if IP != nil && len(IP) == 4 && !IP.IsLoopback() {
        return IP.String(), nil
      }
    }
  }
  return "", errors.New("Could not find interface to bind to")
}

