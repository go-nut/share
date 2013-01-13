package main

import (
  "os"
  "io"
  "archive/tar"
  "compress/gzip"
)

// Write archive of files and dirs to stream w
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
      if err := tarWrite(path, tw); err != nil {
        return err
      }
    }
  }
  return nil
}

// Walk through dir getting files to archive
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
      if err := tarWrite(curPath, tw); err != nil {
        return err
      }
    }
  }
  return nil
}

// Add file to archive stream
func tarWrite(path string, tw *tar.Writer) error {
  file, err := os.Open(path)
  if err != nil {
    return err
  }
  defer file.Close()

  fi, err := file.Stat()
  if err != nil {
    return err
  }

  h := new(tar.Header)
  h.Name = path
  h.Size = fi.Size()
  h.Mode = int64(fi.Mode())
  h.ModTime = fi.ModTime()

  if err := tw.WriteHeader(h); err != nil { 
    return nil
  }
  a, err := io.Copy(tw, file)

  if fi.Size() != a {
    shareLog.Printf("%s did not complete", fi.Name())
  }
  if err != nil {
    shareLog.Print("Error here: ", err.Error())
    return err
  }
  return nil
}

