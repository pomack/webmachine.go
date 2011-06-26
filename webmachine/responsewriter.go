package webmachine

import (
  "http"
  "log"
  "os"
)

func NewResponseWriter(rw http.ResponseWriter) *responseWriter {
  return &responseWriter{rw:rw}
}

func (p *responseWriter) WriteHeader(status int) {
  log.Print("[RW]: Writing Header ", status, "\n")
  p.rw.WriteHeader(status)
}

func (p *responseWriter) Header() http.Header {
  return p.rw.Header()
}

func (p *responseWriter) Write(data []byte) (int, os.Error) {
  log.Print("[RW]: Writing data ", len(data), " bytes\n")
  return p.rw.Write(data)
}
