package main

import (
  "io"
  "os"
  "log"
  "flag"
  "net/http"
  "net"
)

var (
  ip string
  port string
  count int
  self bool
  requests int
  files []string
  shareLog *log.Logger
  wrapper net.Listener
)

func init() {
  flag.StringVar(&port, "p", "8080", "Port to serve on")
  flag.IntVar(&count, "c", 1, "Times to serve")
  flag.BoolVar(&self, "s", false, "Serve self")
}

func main() {

  shareLog = log.New(os.Stdout, "Share ", 1)
  _ip, err := getAddr()
  if err != nil {
    shareLog.Fatal(err)
  }

  flag.StringVar(&ip, "ip", _ip, "IP to serve from")

  flag.Parse()

  files = flag.Args()
  if self {
    files = append(files, os.Args[0])
  }
  if len(files) == 0 {
    files = append(files, ".")
  }

  host := ip + ":" + port

  // Redirect
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if requests < count {
      http.Redirect(w, r, "/download.tar.gz", http.StatusSeeOther)
    }
  })

  http.HandleFunc("/download.tar.gz", dlHandler)

  connChan := make(chan bool, 5)
  wait := make(chan bool)
  go connectionCounter(wait, connChan)

  wrapper, err = wrapListener(host, wait, connChan)
  if err != nil {
    shareLog.Printf("Rrror when attempting to listen: %s", err.Error())
  }
  shareLog.Print("Serving at: http://" + host + "/download.tar.gz")

  err = http.Serve(wrapper, nil)
  <-wait
  shareLog.Print(err.Error())
}

// Handle the download
func dlHandler(w http.ResponseWriter, r *http.Request) {
  requests++
  if requests >= count {
    wrapper.Close()
  }
  shareLog.Printf("Serving to: %s", r.Host)
  w.Header().Set("Content-Type", "application/octet-stream")
  w.Header().Set("Connection", "close")
  w.WriteHeader(http.StatusOK)

  if err := writeArchive(files, w.(io.Writer)); err != nil {
    shareLog.Print("Error serving to: " + r.Host + " - " + err.Error())
  }
  w.(http.Flusher).Flush()
  shareLog.Printf("Finished serving to:  %s", r.Host)
}


// Make sure we are not exiting main before connections are closed
func connectionCounter(wait, connChan chan bool) {
  counter := 0
  keepAlive := true
  var connEvent bool
  for keepAlive || counter != 0 {
    select {
    case connEvent = <-connChan:
      if connEvent {
        counter++
      } else {
        counter--
      }
    case <-wait:
      keepAlive = false
    }
  }
  wait <- true
  shareLog.Print("exiting goroutine")
}
