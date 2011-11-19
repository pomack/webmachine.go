package webmachine

import (
  "container/vector"
  "http"
  "url"
  "io"
  "os"
  "time"
)

type Flusher interface {
  Flush() os.Error
}

type Request interface {
  Method() string  // GET, POST, PUT, etc.
  RawURL() string  // The raw URL given in the request
  URL() *url.URL  // Parsed URL
  Proto() string   // "HTTP/1.0"
  ProtoMajor() int // 1
  ProtoMinor() int // 0
  
  Header() http.Header
  Cookie(name string) (*http.Cookie, os.Error)
  Cookies() []*http.Cookie
  Body() io.ReadCloser
  ContentLength() int64
  TransferEncoding() []string
  CloseAfterReply() bool
  Host() string
  Referer() string
  UserAgent() string
  Form() map[string][]string
  Trailer() http.Header
  
  HostParts() []string
  URLParts() []string
}

type Context interface {}

type request struct {
  req *http.Request
  hostParts []string
  urlParts []string
}

type RouteHandler interface {
  HandlerFor(Request, ResponseWriter) RequestHandler
}

type MediaTypeHandler interface {
  MediaType() string
  OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter)
}

type MediaTypeInputHandler interface {
  MediaType() string
  OutputTo(req Request, cxt Context, writer io.Writer) (int, http.Header, os.Error)
}

type CharsetHandler interface {
  Charset() string
  CharsetConverter(req Request, cxt Context, reader io.Reader) (io.Reader)
}

type EncodingHandler interface {
  Encoding() string
  Encoder(req Request, cxt Context, writer io.Writer) (io.Writer)
  Decoder(req Request, cxt Context, reader io.Reader) (io.Reader)
}

type RequestHandler interface {
  StartRequest(req Request, cxt Context) (Request, Context)
  ResourceExists(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  ServiceAvailable(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  IsAuthorized(req Request, cxt Context) (bool, string, Request, Context, int, os.Error)
  Forbidden(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  AllowMissingPost(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  MalformedRequest(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  URITooLong(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  KnownContentType(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  ValidContentHeaders(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  ValidEntityLength(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  Options(req Request, cxt Context) ([]string, Request, Context, int, os.Error)
  AllowedMethods(req Request, cxt Context) ([]string, Request, Context, int, os.Error)
  DeleteResource(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  DeleteCompleted(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  PostIsCreate(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  CreatePath(req Request, cxt Context) (string, Request, Context, int, os.Error)
  ProcessPost(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  ContentTypesProvided(req Request, cxt Context) ([]MediaTypeHandler, Request, Context, int, os.Error)
  ContentTypesAccepted(req Request, cxt Context) ([]MediaTypeInputHandler, Request, Context, int, os.Error)
  IsLanguageAvailable(languages []string, req Request, cxt Context) (bool, Request, Context, int, os.Error)
  CharsetsProvided(charsets []string, req Request, cxt Context) ([]CharsetHandler, Request, Context, int, os.Error)
  EncodingsProvided(encodings []string, req Request, cxt Context) ([]EncodingHandler, Request, Context, int, os.Error)
  Variances(req Request, cxt Context) ([]string, Request, Context, int, os.Error)
  IsConflict(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  MultipleChoices(req Request, cxt Context) (bool, http.Header, Request, Context, int, os.Error)
  PreviouslyExisted(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  MovedPermanently(req Request, cxt Context) (string, Request, Context, int, os.Error)
  MovedTemporarily(req Request, cxt Context) (string, Request, Context, int, os.Error)
  LastModified(req Request, cxt Context) (*time.Time, Request, Context, int, os.Error)
  Expires(req Request, cxt Context) (*time.Time, Request, Context, int, os.Error)
  GenerateETag(req Request, cxt Context) (string, Request, Context, int, os.Error)
  FinishRequest(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  ResponseIsRedirect(req Request, cxt Context) (bool, Request, Context, int, os.Error)
  
  HasRespBody(req Request, cxt Context) bool
}

type WebMachine interface {
  ServeHTTP(http.ResponseWriter, *http.Request)
  AddRouteHandler(RouteHandler)
  RemoveRouteHandler(RouteHandler)
}

type webMachine struct {
  routeHandlers vector.Vector;
}

type WriteThrough struct {
  from io.Writer
  to io.Writer
}

