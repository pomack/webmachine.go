package webmachine

import (
    "io"
    "log"
    "net/http"
)

type ResponseWriter interface {
    http.ResponseWriter
    io.Closer
    Flusher
    AddEncoding(h EncodingHandler, req Request, cxt Context) io.Writer
}

type responseWriter struct {
    rw  http.ResponseWriter
    w   io.Writer
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
    return &responseWriter{rw: rw, w: rw}
}

func (p *responseWriter) WriteHeader(status int) {
    log.Print("[RW]: Writing Header ", status)
    p.rw.WriteHeader(status)
}

func (p *responseWriter) Header() http.Header {
    return p.rw.Header()
}

func (p *responseWriter) Write(data []byte) (int, error) {
    log.Print("[RW]: Writing data ", len(data), " bytes")
    if len(data) < 5000 {
        log.Print("[RW]: Wrote:\n", string(data))
    }
    return p.w.Write(data)
}

func (p *responseWriter) AddEncoding(h EncodingHandler, req Request, cxt Context) io.Writer {
    writer := h.Encoder(req, cxt, p.w)
    if writer != nil {
        p.w = writer
    }
    return p.w
}

func (p *responseWriter) Flush() error {
    if p.rw != p.w {
        if f, ok := p.w.(Flusher); ok {
            log.Print("[RW]: Flushing Writer")
            return f.Flush()
        }
    }
    if f, ok := p.rw.(Flusher); ok {
        log.Print("[RW]: Flushing ResponseWriter")
        return f.Flush()
    }
    log.Print("[RW]: Failed to Flush, not flushable")
    return nil
}

func (p *responseWriter) Close() error {
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
