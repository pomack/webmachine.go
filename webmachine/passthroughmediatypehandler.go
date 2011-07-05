package webmachine

import (
  "http"
  "io"
  "json"
  "log"
  "path"
  "strconv"
  "time"
  "os"
)


type PassThroughMediaTypeHandler struct {
  mediaType string
  reader io.ReadCloser
  numberOfBytes int64
  lastModified *time.Time
  writtenStatusHeader bool
}


type PassThroughMediaTypeInputHandler struct {
  mediaType string
  charset string
  language string
  filename string
  urlPath string
  append bool
  numberOfBytes int64
  reader io.Reader
  writtenStatusHeader bool
}




func NewPassThroughMediaTypeHandler(mediaType string, reader io.ReadCloser, numberOfBytes int64, lastModified *time.Time) *PassThroughMediaTypeHandler {
  return &PassThroughMediaTypeHandler{
    mediaType: mediaType,
    reader: reader,
    numberOfBytes: numberOfBytes,
    lastModified: lastModified,
  }
}

func (p *PassThroughMediaTypeHandler) MediaType() string {
  return p.mediaType
}

func (p *PassThroughMediaTypeHandler) OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
  if !p.writtenStatusHeader {
    resp.WriteHeader(200)
    p.writtenStatusHeader = true
  }
  if req.Header().Get("Accept-Ranges") == "bytes" {
    rangeHeader := req.Header().Get("Range")
    if len(rangeHeader) > 6 && rangeHeader[0:6] == "bytes=" {
      ranges := p.splitRangeHeaderString(rangeHeader)
      outRangeString := "bytes="
      for i, arange := range ranges {
        if i > 0 {
          outRangeString += ","
        }
        outRangeString += strconv.Itoa64(arange[0]) + "-" + strconv.Itoa64(arange[1]-1)
      }
      outRangeString += "/" + strconv.Itoa64(p.numberOfBytes)
      resp.Header().Set("Content-Range", "bytes=" + outRangeString)
      currentOffset := int64(0)
      for _, arange := range ranges {
        start := arange[0]
        end := arange[1]
        var err os.Error
        if currentOffset < start {
          if seeker, ok := p.reader.(io.Seeker); ok {
            currentOffset, err = seeker.Seek(start - currentOffset, 1)
            if err != nil {
              return
            }
          } else {
            if start - currentOffset >= 32768 {
              buf := make([]byte, 32768)
              for ; currentOffset + 32768 < start; currentOffset += 32768 {
                if _, err = io.ReadFull(p.reader, buf); err != nil {
                  return
                }
              }
            }
            if currentOffset < start {
              buf := make([]byte, start - currentOffset)
              if _, err = io.ReadFull(p.reader, buf); err != nil {
                return
              }
            }
          }
        }
        if req.Method() == HEAD {
          return
        }
        for currentOffset < end {
          written, err := io.Copyn(writer, p.reader, end - currentOffset)
          currentOffset += written
          if err != nil {
            return
          }
        }
      }
      return
    }
  }
  if req.Method() == HEAD {
    return
  }
  currentOffset := int64(0)
  log.Print("[PTMTH]: Writer ", writer, "\n[PTMTH]: Reader ", p.reader, "\n[PTMTH]: numBytes ", p.numberOfBytes, "\n[PTMTH]: currentOffset ", currentOffset, "\n")
  for currentOffset < int64(p.numberOfBytes) {
    bytesToSend := p.numberOfBytes - currentOffset
    data := make([]byte, bytesToSend  + 10000)
    numBytesRead, err := p.reader.Read(data[0:bytesToSend])
    currentOffset += int64(numBytesRead)
    if err != nil {
      return
    }
    log.Print("[PTMTH]: About to write ", len(data[0:bytesToSend]), " bytes to the writer\n")
    _, err = writer.Write(data[0:bytesToSend])
    if err != nil {
      return
    }
    //written, err := io.Copyn(writer, p.reader, p.numberOfBytes - currentOffset)
    //if err != nil {
    //  return
    //}
    //currentOffset += int64(written)
  }
}






func NewPassThroughMediaTypeInputHandler(mediaType, charset, language, filename, urlPath string, append bool, numberOfBytes int64, reader io.Reader) *PassThroughMediaTypeInputHandler {
  return &PassThroughMediaTypeInputHandler{
    mediaType: mediaType,
    charset: charset,
    language: language,
    filename: filename,
    urlPath: urlPath,
    append: append,
    numberOfBytes: numberOfBytes,
    reader: reader,
  }
}

func (p *PassThroughMediaTypeInputHandler) MediaType() string {
  return p.mediaType
}

func (p *PassThroughMediaTypeInputHandler) OutputTo(req Request, cxt Context, writer io.Writer) (int, http.Header, os.Error) {
  fileInfo, err := os.Stat(p.filename)
  var file *os.File
  m := make(map[string]string)
  w := json.NewEncoder(writer)
  dirname, _ := path.Split(p.filename)
  file = nil
  defer func() {
    if file != nil {
      file.Close()
    }
  }()
  if fileInfo == nil {
    if err = os.MkdirAll(dirname, 0644); err != nil {
      log.Print("[PTMTIH]: Unable to create directory to store file due to error: ", err)
      headers := make(http.Header)
      headers.Set("Content-Type", "application/json")
      m["status"] = "error"
      m["message"] = err.String()
      m["result"] = p.urlPath
      w.Encode(m)
      return 500, headers, err
    }
    if file, err = os.OpenFile(p.filename, os.O_CREATE, 0644); err != nil {
      log.Print("[PTMTIH]: Unable to create file named: \"", p.filename, "\" due to error: ", err)
      headers := make(http.Header)
      headers.Set("Content-Type", "application/json")
      m["status"] = "error"
      m["message"] = err.String()
      m["result"] = p.urlPath
      w.Encode(m)
      return 500, headers, err
    }
  } else {
    if p.append {
      file, err = os.OpenFile(p.filename, os.O_APPEND, 0644)
    } else {
      file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_TRUNC, 0644)
    }
    if err != nil {
      log.Print("[PTMTIH]: Unable to open file \"", p.filename, "\"for writing due to error: ", err)
      headers := make(http.Header)
      headers.Set("Content-Type", "application/json")
      m["status"] = "error"
      m["message"] = err.String()
      m["result"] = p.urlPath
      w.Encode(m)
      return 500, headers, err
    }
  }
  var n int64
  if p.numberOfBytes >= 0 {
    n, err = io.Copyn(file, p.reader, p.numberOfBytes)
  } else {
    n, err = io.Copy(file, p.reader)
  }
  log.Print("[PTMTIH]: Wrote ", n, " bytes to file with error: ", err)
  if err != nil && err != os.EOF {
    headers := make(http.Header)
    headers.Set("Content-Type", "application/json")
    m["status"] = "error"
    m["message"] = err.String()
    m["result"] = p.urlPath
    w.Encode(m)
    return 500, headers, err
  }
  headers := make(http.Header)
  headers.Set("Content-Type", "application/json")
  m["status"] = "success"
  m["message"] = ""
  m["result"] = p.urlPath
  w.Encode(m)
  return 200, headers, nil
}
