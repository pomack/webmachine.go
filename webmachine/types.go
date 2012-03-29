package webmachine

import (
    "io"
    "mime/multipart"
    "net/http"
    "net/url"
    "time"
)

type Flusher interface {
    Flush() error
}

type Request interface {
    UnderlyingRequest() *http.Request
    Method() string  // GET, POST, PUT, etc.
    RawURL() string  // The raw URL given in the request
    URL() *url.URL   // Parsed URL
    Proto() string   // "HTTP/1.0"
    ProtoMajor() int // 1
    ProtoMinor() int // 0

    Header() http.Header
    AddCookie(c *http.Cookie)
    Cookie(name string) (*http.Cookie, error)
    Cookies() []*http.Cookie
    Body() io.ReadCloser
    ContentLength() int64
    TransferEncoding() []string
    CloseAfterReply() bool
    Host() string
    Referer() string
    UserAgent() string
    Form() map[string][]string
    FormFile(key string) (multipart.File, *multipart.FileHeader, error)
    FormValue(key string) string
    MultipartReader() (*multipart.Reader, error)
    ParseForm() (err error)
    ParseMultipartForm(maxMemory int64) error
    Trailer() http.Header
    HostParts() []string
    URLParts() []string
}

type Context interface{}

type request struct {
    req       *http.Request
    hostParts []string
    urlParts  []string
}

type RouteHandler interface {
    HandlerFor(Request, ResponseWriter) RequestHandler
}

type MediaTypeHandler interface {
    MediaTypeOutput() string
    MediaTypeHandleOutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter)
}

type MediaTypeInputHandler interface {
    MediaTypeInput() string
    MediaTypeHandleInputFrom(req Request, cxt Context) (int, http.Header, io.WriterTo)
}

type CharsetHandler interface {
    Charset() string
    CharsetConverter(req Request, cxt Context, reader io.Reader) io.Reader
}

type EncodingHandler interface {
    Encoding() string
    Encoder(req Request, cxt Context, writer io.Writer) io.Writer
    Decoder(req Request, cxt Context, reader io.Reader) io.Reader
}

type RequestHandler interface {
    StartRequest(req Request, cxt Context) (Request, Context)
    ResourceExists(req Request, cxt Context) (bool, Request, Context, int, error)
    ServiceAvailable(req Request, cxt Context) (bool, Request, Context, int, error)
    IsAuthorized(req Request, cxt Context) (bool, string, Request, Context, int, error)
    Forbidden(req Request, cxt Context) (bool, Request, Context, int, error)
    AllowMissingPost(req Request, cxt Context) (bool, Request, Context, int, error)
    MalformedRequest(req Request, cxt Context) (bool, Request, Context, int, error)
    URITooLong(req Request, cxt Context) (bool, Request, Context, int, error)
    KnownContentType(req Request, cxt Context) (bool, Request, Context, int, error)
    ValidContentHeaders(req Request, cxt Context) (bool, Request, Context, int, error)
    ValidEntityLength(req Request, cxt Context) (bool, Request, Context, int, error)
    Options(req Request, cxt Context) ([]string, Request, Context, int, error)
    AllowedMethods(req Request, cxt Context) ([]string, Request, Context, int, error)
    DeleteResource(req Request, cxt Context) (bool, Request, Context, int, error)
    DeleteCompleted(req Request, cxt Context) (bool, Request, Context, int, error)
    PostIsCreate(req Request, cxt Context) (bool, Request, Context, int, error)
    CreatePath(req Request, cxt Context) (string, Request, Context, int, error)
    ProcessPost(req Request, cxt Context) (Request, Context, int, http.Header, io.WriterTo, error)
    ContentTypesProvided(req Request, cxt Context) ([]MediaTypeHandler, Request, Context, int, error)
    ContentTypesAccepted(req Request, cxt Context) ([]MediaTypeInputHandler, Request, Context, int, error)
    IsLanguageAvailable(languages []string, req Request, cxt Context) (bool, Request, Context, int, error)
    CharsetsProvided(charsets []string, req Request, cxt Context) ([]CharsetHandler, Request, Context, int, error)
    EncodingsProvided(encodings []string, req Request, cxt Context) ([]EncodingHandler, Request, Context, int, error)
    Variances(req Request, cxt Context) ([]string, Request, Context, int, error)
    IsConflict(req Request, cxt Context) (bool, Request, Context, int, error)
    MultipleChoices(req Request, cxt Context) (bool, http.Header, Request, Context, int, error)
    PreviouslyExisted(req Request, cxt Context) (bool, Request, Context, int, error)
    MovedPermanently(req Request, cxt Context) (string, Request, Context, int, error)
    MovedTemporarily(req Request, cxt Context) (string, Request, Context, int, error)
    LastModified(req Request, cxt Context) (time.Time, Request, Context, int, error)
    Expires(req Request, cxt Context) (time.Time, Request, Context, int, error)
    GenerateETag(req Request, cxt Context) (string, Request, Context, int, error)
    FinishRequest(req Request, cxt Context) (bool, Request, Context, int, error)
    ResponseIsRedirect(req Request, cxt Context) (bool, Request, Context, int, error)

    HasRespBody(req Request, cxt Context) bool
}

type WebMachine interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
    AddRouteHandler(RouteHandler)
    RemoveRouteHandler(RouteHandler)
}

type webMachine struct {
    routeHandlers []RouteHandler
}

type WriteThrough struct {
    from io.Writer
    to   io.Writer
}
