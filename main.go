package main


import (
  "html"
  "flag"
  "net"
  "net/http"
  "archive/tar"
  "compress/gzip"
  "log"
  "os"
  "io"
)

var (
  ip string
  port string
  count int
  requests int
//  self bool
  errLog *log.Logger
  outLog *log.Logger
  l net.Listener
  fileList []string
)

func init() {
  flag.StringVar(&port, "p", "8080", "Port to serve on")
  flag.IntVar(&count, "c", 1, "Times to serve")
//  flag.BoolVar(&self, "s", false, "Serve self")
  errLog = log.New(os.Stderr, "Error ", 2)
  outLog = log.New(os.Stdout, "Share: ", 1)


  requests = 0
}

func main() {

//  a, err := net.InterfaceAddrs()
//  logErr(err)
//  _ip, _, err := net.ParseCIDR(a[len(a)-1].String())

  var _ip string
  interfaces, err := net.Interfaces()
  if err != nil {
    errLog.Println(err)
  }
  for _, interf := range interfaces {
    addrs, err := interf.Addrs()
    if err != nil {
      errLog.Println(err)
    }
    for _, addr := range addrs {
      IP, _, err := net.ParseCIDR(addr.String())
      if err != nil {
        errLog.Println(err)
        continue
      }
      if len(IP) == 4 && !IP.IsLoopback() {
        _ip = IP.String()
        break
      }
    }
  }
  flag.StringVar(&ip, "ip", _ip, "IP to serve from")

  flag.Parse()

  host := ip + ":" + port
  outLog.Println("Serving on http://" + host + "/download.tar.gz")

  fileList = flag.Args()
  if len(fileList) == 0 {
    fileList = append(fileList, ".")
  }

  laddr, err := net.ResolveTCPAddr("tcp", host)
  if err != nil {
    errLog.Fatalln(err)
  }

  l, err = net.ListenTCP("tcp", laddr)
  if err != nil {
    errLog.Fatalln(err)
  }

  http.HandleFunc("/", handler)
  err = http.Serve(l, nil)
  if err != nil {
    errLog.Fatalln(err)
  }
}

func logErr(err error) {
  if err != nil {
    errLog.Fatalln(err)
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
  outLog.Println("starting handler")
  if requests >= count {
    return
  }

  if q := html.EscapeString(r.URL.Path); q != "/download.tar.gz" {
    http.Redirect(w, r, "download.tar.gz", 302)
    return
  }

  w.Header().Set("Content-Type", "application/octet-stream")
  writer := w.(io.Writer)
  outLog.Println("Serving to " + r.Host)
  if err := writeArchive(fileList, writer); err != nil {
    errLog.Fatalln(err)
  }
  outLog.Println("Served to " + r.Host + " complete")

  requests++
  if requests >= count {
    defer closeConn()
  }
  outLog.Println("ending handler")
}

func iterDir(dirPath string, tw *tar.Writer) error {
  dir, err := os.Open( dirPath )
  if err != nil {
    return err
  }
  defer dir.Close()

  fis, err := dir.Readdir( 0 )
  if err != nil {
    return err
  }

  for _, fi := range fis {
    curPath := dirPath + "/" + fi.Name()
    if fi.IsDir() {
      if err := iterDir(curPath, tw); err != nil {
        return err
      }
    } else {
      if err := tarWrite(curPath, tw, fi); err != nil {
        return err
      }
    }
  }
  return nil
}

func tarWrite(path string, tw *tar.Writer, fi os.FileInfo) error {
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  outLog.Println(path)
  h := new(tar.Header)
  h.Name = path
  h.Size = fi.Size()
  h.Mode = int64(fi.Mode())
  h.ModTime = fi.ModTime()

  err = tw.WriteHeader(h)
  if err != nil {
    return err
  }
  _, err = io.Copy(tw, file)
  if err != nil {
    return err
  }
  return nil
}

func writeArchive(paths []string, w io.Writer) error {
  gz := gzip.NewWriter(w)
  defer gz.Close()

  tw := tar.NewWriter(gz)
  defer tw.Close()

  for _, path := range paths {
    f, err := os.Open(path)
    if err != nil {
      return err
    }

    stat, err := f.Stat()
    if err != nil {
      return err
    }

    if stat.IsDir() {
      if err := iterDir(path, tw); err != nil {
        return err
      }
    } else {
      if err := tarWrite(path, tw, stat); err != nil {
        return err
      }
    }
  }
  return nil
}

func closeConn() {
  if err := l.Close(); err != nil {
    errLog.Fatalln(err)
  }
}
