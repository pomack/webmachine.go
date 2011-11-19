package webmachine

import (
  "http"
  "io"
  "json"
  "log"
  "path"
  "os"
)


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
      headers.Set("Content-Type", MIME_TYPE_JSON)
      m["status"] = "error"
      m["message"] = err.String()
      m["result"] = p.urlPath
      w.Encode(m)
      return 500, headers, err
    }
    if file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
      log.Print("[PTMTIH]: Unable to create file named: \"", p.filename, "\" due to error: ", err)
      headers := make(http.Header)
      headers.Set("Content-Type", MIME_TYPE_JSON)
      m["status"] = "error"
      m["message"] = err.String()
      m["result"] = p.urlPath
      w.Encode(m)
      return 500, headers, err
    }
  } else {
    if p.append {
      file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_APPEND, 0644)
    } else {
      file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_TRUNC, 0644)
    }
    if err != nil {
      log.Print("[PTMTIH]: Unable to open file \"", p.filename, "\"for writing due to error: ", err)
      headers := make(http.Header)
      headers.Set("Content-Type", MIME_TYPE_JSON)
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
    headers.Set("Content-Type", MIME_TYPE_JSON)
    m["status"] = "error"
    m["message"] = err.String()
    m["result"] = p.urlPath
    w.Encode(m)
    return 500, headers, err
  }
  headers := make(http.Header)
  headers.Set("Content-Type", MIME_TYPE_JSON)
  m["status"] = "success"
  m["message"] = ""
  m["result"] = p.urlPath
  w.Encode(m)
  return 200, headers, nil
}
