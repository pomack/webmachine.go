package webmachine

import (
    "io"
    "mime/multipart"
    "net/http"
    "net/url"
    "strings"
)

func NewRequestFromHttpRequest(req *http.Request) Request {
    p := new(request)
    p.req = req
    p.hostParts = strings.Split(req.Host, ".")
    p.urlParts = strings.Split(req.URL.Path, "/")
    return p
}

func (p *request) UnderlyingRequest() *http.Request {
    return p.req
}

func (p *request) Method() string {
    return p.req.Method
}

func (p *request) RawURL() string {
    return p.req.URL.String()
}

func (p *request) URL() *url.URL {
    return p.req.URL
}

func (p *request) Proto() string {
    return p.req.Proto
}

func (p *request) ProtoMajor() int {
    return p.req.ProtoMajor
}

func (p *request) ProtoMinor() int {
    return p.req.ProtoMinor
}

func (p *request) Header() http.Header {
    return p.req.Header
}

func (p *request) AddCookie(c *http.Cookie) {
    p.req.AddCookie(c)
}

func (p *request) Cookie(name string) (*http.Cookie, error) {
    return p.req.Cookie(name)
}

func (p *request) Cookies() []*http.Cookie {
    return p.req.Cookies()
}

func (p *request) Body() io.ReadCloser {
    return p.req.Body
}

func (p *request) ContentLength() int64 {
    return p.req.ContentLength
}

func (p *request) TransferEncoding() []string {
    return p.req.TransferEncoding
}

func (p *request) CloseAfterReply() bool {
    return p.req.Close
}

func (p *request) Host() string {
    return p.req.Host
}

func (p *request) Referer() string {
    return p.req.Referer()
}

func (p *request) UserAgent() string {
    return p.req.UserAgent()
}

func (p *request) Form() map[string][]string {
    return p.req.Form
}

func (p *request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
    return p.req.FormFile(key)
}

func (p *request) FormValue(key string) string {
    return p.req.FormValue(key)
}

func (p *request) MultipartReader() (*multipart.Reader, error) {
    return p.req.MultipartReader()
}

func (p *request) ParseForm() error {
    return p.req.ParseForm()
}

func (p *request) ParseMultipartForm(maxMemory int64) error {
    return p.req.ParseMultipartForm(maxMemory)
}

func (p *request) Trailer() http.Header {
    return p.req.Trailer
}

func (p *request) HostParts() []string {
    return p.hostParts
}

func (p *request) URLParts() []string {
    return p.urlParts
}
