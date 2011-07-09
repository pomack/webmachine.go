package webmachine

import (
  "http"
  "log"
  "io"
  "os"
)

type ResponseWriter interface {
  http.ResponseWriter
  io.Closer
  Flusher
  AddEncoding(h EncodingHandler, req Request, cxt Context) (io.Writer)
}


type responseWriter struct {
  rw http.ResponseWriter
  w io.Writer
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
  return &responseWriter{rw:rw, w:rw}
}

func (p *responseWriter) WriteHeader(status int) {
  log.Print("[RW]: Writing Header ", status)
  p.rw.WriteHeader(status)
}

func (p *responseWriter) Header() http.Header {
  return p.rw.Header()
}

func (p *responseWriter) Write(data []byte) (int, os.Error) {
  log.Print("[RW]: Writing data ", len(data), " bytes")
  return p.w.Write(data)
}

func (p *responseWriter) AddEncoding(h EncodingHandler, req Request, cxt Context) (io.Writer) {
  writer := h.Encoder(req, cxt, p.w)
  if writer != nil {
    p.w = writer
  }
  return p.w
}

func (p *responseWriter) Flush() (os.Error) {
  if p.rw != p.w {
    if f, ok := p.w.(Flusher); ok {
      log.Print("[RW]: Flushing Writer")
      return f.Flush()
    }
    if c, ok := p.w.(io.Closer); ok {
      log.Print("[RW]: Closing Writer on Flush")
      return c.Close()
    }
  }
  if f, ok := p.rw.(Flusher); ok {
    log.Print("[RW]: Flushing ResponseWriter")
    return f.Flush()
  }
  if c, ok := p.rw.(io.Closer); ok {
    log.Print("[RW]: Closing ResponseWriter")
    return c.Close()
  }
  log.Print("[RW]: Failed to Flush, not flushable")
  return nil
}

func (p *responseWriter) Close() (os.Error) {
  if p.rw != p.w {
    if closer, ok := p.w.(io.Closer); ok {
      closer.Close()
    }
  }
  if closer, ok := p.rw.(io.Closer); ok {
    closer.Close()
  }
  return nil
}
