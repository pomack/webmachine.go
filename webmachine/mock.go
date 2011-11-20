package webmachine

import (
    "bytes"
    "http"
    "io/ioutil"
    "json"
    "os"
)

type MockResponseWriter struct {
    Headers         http.Header     `json:"headers,omitempty"`
    Buffer          *bytes.Buffer   `json:"buffer,omitempty"`
    StatusCode      int             `json:"status_code,omitempty"`
    Request         *http.Request   `json:"request,omitempty"`
}

func NewMockResponseWriter(request *http.Request) http.ResponseWriter {
    return &MockResponseWriter{
        Headers: make(http.Header),
        Buffer: bytes.NewBufferString(""),
        StatusCode: 0,
        Request: request,
    }
}

func (p *MockResponseWriter) Header() http.Header {
    return p.Headers
}

func (p *MockResponseWriter) Write(data []byte) (int, os.Error) {
    return p.Buffer.Write(data)
}

func (p *MockResponseWriter) WriteHeader(statusCode int) {
    p.StatusCode = statusCode
}

func (p *MockResponseWriter) MarshalJSON() ([]byte, os.Error) {
    m := make(map[string]interface{})
    m["headers"] = p.Headers
    m["buffer"] = p.Buffer.String()
    m["status_code"] = p.StatusCode
    return json.Marshal(m)
}

func (p *MockResponseWriter) String() string {
    resp := new(http.Response)
    resp.StatusCode = p.StatusCode
    resp.Proto = "HTTP/1.1"
    resp.ProtoMajor = 1
    resp.ProtoMinor = 1
    resp.Header = p.Headers
    resp.Body = ioutil.NopCloser(bytes.NewBuffer(p.Buffer.Bytes()))
    resp.ContentLength = int64(p.Buffer.Len())
    if p.Headers.Get("Transfer-Encoding") != "" {
        resp.TransferEncoding = []string{p.Headers.Get("Transfer-Encoding")}
    } else {
        resp.TransferEncoding = nil
    }
    resp.Close = p.Headers.Get("Connection") == "close"
    resp.Trailer = make(http.Header)
    resp.Request = p.Request
    b, _ := http.DumpResponse(resp, true)
    return string(b)
}

