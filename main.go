package main

import (
  "io"
  "os"
  "log"
  "flag"
  "net/http"
)

var (
  ip string
  port string
  count int
  requests int
  files []string
  httpServer *StopServer
  shareLog *log.Logger
)

func init() {
  flag.StringVar(&port, "p", "8080", "Port to serve on")
  flag.IntVar(&count, "c", 1, "Times to serve")
//flag.BoolVar(&self, "s", false, "Serve self")
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
  if len(files) == 0 {
    files = append(files, ".")
  }

  host := ip + ":" + port

  httpServer = &StopServer {
    http.Server: http.Server {
      Addr: host,
    },
  }

  handler := http.NewServeMux()

  // Redirect
  handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if requests < count {
      http.Redirect(w, r, "/download.tar.gz", http.StatusSeeOther)
    }
  })

  handler.HandleFunc("/download.tar.gz", dlHandler)

  httpServer.Handler = handler

  shareLog.Print("Serving at: http://" + host + "/download.tar.gz")
  err = httpServer.ListenAndServe()

  httpServer.WaitUnfinished()
}

// Handle the download
func dlHandler(w http.ResponseWriter, r *http.Request) {
  requests++
  if requests >= count {
    httpServer.Stop()
  }
  w.Header().Set("Content-Type", "application/octet-stream")
  writer := w.(io.Writer)
  shareLog.Print("Serving to: " + r.Host)
  if err := writeArchive(files, writer); err != nil {
    shareLog.Print("Error serving to: " + r.Host + " - " + err.Error())
  }
  shareLog.Print("Finished serving to: " + r.Host)
}
