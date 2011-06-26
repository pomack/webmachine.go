package webmachine

import (
  "template"
)

const (
  CONNECT = "CONNECT"
  DELETE = "DELETE"
  GET = "GET"
  HEAD = "HEAD"
  OPTIONS = "OPTIONS"
  PUT = "PUT"
  POST = "POST"
  TRACE = "TRACE"
)

const (
  ISO_8601_DATETIME_FORMAT = "2006-01-02T03:04:05Z"
)

var ALL_METHODS []string
var HTML_DIRECTORY_LISTING_ERROR_TEMPLATE *template.Template
var HTML_DIRECTORY_LISTING_SUCCESS_TEMPLATE *template.Template

type WMDecision int

const (
  wmResponded WMDecision = iota
  v3b13 // Service Available?
  v3b13b 
  v3b12 // Known method?
  v3b11 // URI too long?
  v3b10 // Method allowed?
  v3b9 // Malformed?
  v3b8 // Authorized?
  v3b7 // Forbidden?
  v3b6 // Okay Content-* Headers?
  v3b5 // Known Content-Type?
  v3b4 // Req Entity Too Large?
  v3b3 // OPTIONS?
  v3c3 // Accept exists?
  v3c4 // Acceptable media type available?
  v3d4 // Accept-Language exists?
  v3d5 // Acceptable Language available?
  v3e5 // Accept-Charset exists?
  v3e6 // Accceptable Charset available?
  v3f6 // Accept-Encoding exists?
  v3f7 // Acceptable encoding available?
  v3g7 // Resource exists?
  v3g8 // If-Match exists?
  v3g9 // If-Match: * exists
  v3g11 // ETag in If-Match
  v3h7 // If-Match: * exists
  v3h10 // If-unmodified-since exists?
  v3h11 // I-UM-S is valid date?
  v3h12 // Last-Modified > I-UM-S?
  v3i4 // Moved permanently?
  v3i7 // PUT?
  v3i12 // If-none-match exists?
  v3i13 // If-None-Match: * exists?
  v3j18 // GET or HEAD?
  v3k5 // Moved permanently?
  v3k7 // Previously existed?
  v3k13 // Etag in if-none-match?
  v3l5 // Moved temporarily?
  v3l7 // POST?
  v3l13 // IMS exists?
  v3l14 // IMS is valid date?
  v3l15 // IMS > Now?
  v3l17 // Last-Modified > IMS?
  v3m5 // POST?
  v3m7 // Server allows POST to missing resource?
  v3m16 // DELETE?
  v3m20 // DELETE enacted immediately?
  v3m20b //
  v3n5 // Server allows POST to missing resource?
  v3n11 // Redirect?
  v3n16 // POST?
  v3o14 // Conflict?
  v3o16 // PUT?
  v3o18 // Multiple representations?
  v3o18b //
  v3o20 // Response includes an entity?
  v3p3 // Conflict?
  v3p11 // New resource?
)


var (
  defaultMimeTypes map[string]string
)


func init() {
  ALL_METHODS = []string{GET, HEAD, POST, CONNECT, DELETE, OPTIONS, PUT, TRACE}
  
  defaultMimeTypes = make(map[string]string)
  defaultMimeTypes[".htm"] = "text/html"
  defaultMimeTypes[".html"] = "text/html"
  defaultMimeTypes[".xhtml"] = "application/xhtml+xml"
  defaultMimeTypes[".xml"] = "application/xml"
  defaultMimeTypes[".css"] = "text/css"
  defaultMimeTypes[".js"] = "application/x-javascript"
  defaultMimeTypes[".json"] = "application/json"
  defaultMimeTypes[".jpg"] = "image/jpeg"
  defaultMimeTypes[".jpeg"] = "image/jpeg"
  defaultMimeTypes[".gif"] = "image/gif"
  defaultMimeTypes[".png"] = "image/png"
  defaultMimeTypes[".ico"] = "image/x-icon"
  defaultMimeTypes[".swf"] = "application/x-shockwave-flash"
  defaultMimeTypes[".zip"] = "application/zip"
  defaultMimeTypes[".bz2"] = "application/x-bzip2"
  defaultMimeTypes[".gz"] = "application/x-gzip"
  defaultMimeTypes[".tar"] = "application/x-tar"
  defaultMimeTypes[".tgz"] = "application/x-gzip"
  defaultMimeTypes[".htc"] = "text/x-component"
  defaultMimeTypes[".manifest"] = "text/cache-manifest"
  defaultMimeTypes[".svg"] = "image/svg+xml"
  defaultMimeTypes[".txt"] = "text/plain"
  defaultMimeTypes[".text"] = "text/plain"
  defaultMimeTypes[".csv"] = "text/csv"
  
  HTML_DIRECTORY_LISTING_SUCCESS_TEMPLATE = template.MustParseFile("templates/html/directory_listing/success.html", nil)
  HTML_DIRECTORY_LISTING_ERROR_TEMPLATE = template.MustParseFile("templates/html/directory_listing/error.html", nil)
}

func (p WMDecision) String() string {
  var s string
  switch p {
  case wmResponded:
    s = "Responded"
  case v3b13: // Service Available?
    s = "v3b13: Service Available?"
  case v3b13b:
    s = "v3b13b: Service Available?"
  case v3b12: // Known method?
    s = "v3b12: Known method?"
  case v3b11: // URI too long?
    s = "v3b11: URI too long?"
  case v3b10: // Method allowed?
    s = "v3b10: Method allowed?"
  case v3b9: // Malformed?
    s = "v3b9: Malformed?"
  case v3b8: // Authorized?
    s = "v3b8: Authorized?"
  case v3b7: // Forbidden?
    s = "v3b7: Forbidden?"
  case v3b6: // Okay Content-* Headers?
    s = "v3b6: Okay Content-* Headers?"
  case v3b5: // Known Content-Type?
    s = "v3b5: Known Content-Type?"
  case v3b4: // Req Entity Too Large?
    s = "v3b4: Req Entity Too Large?"
  case v3b3: // OPTIONS?
    s = "v3b3: OPTIONS?"
  case v3c3: // Accept exists?
    s = "v3c3: Accept exists?"
  case v3c4: // Acceptable media type available?
    s = "v3c4: Acceptable media type available?"
  case v3d4: // Accept-Language exists?
    s = "v3d4: Accept-Language exists?"
  case v3d5: // Acceptable Language available?
    s = "v3d5: Acceptable Language available?"
  case v3e5: // Accept-Charset exists?
    s = "v3e5: Accept-Charset exists?"
  case v3e6: // Accceptable Charset available?
    s = "v3e6: Acceptable Charset available?"
  case v3f6: // Accept-Encoding exists?
    s = "v3f6: Accept-Encoding exists?"
  case v3f7: // Acceptable encoding available?
    s = "v3f7: Acceptable encoding available?"
  case v3g7: // Resource exists?
    s = "v3g7: Resource exists?"
  case v3g8: // If-Match exists?
    s = "v3g8: If-Match exists?"
  case v3g9: // ETag in If-Match
    s = "v3g9: ETag in If-Match"
  case v3g11: // ETag in If-Match
    s = "v3g11: ETag in If-Match"
  case v3h7: // If-Match: * exists
    s = "v3h7: If-Match: * exists"
  case v3h10: // If-unmodified-since exists?
    s = "v3h10: If-unmodified-since exists?"
  case v3h11: // I-UM-S is valid date?
    s = "v3h11: I-UM-S is valid date?"
  case v3h12: // Last-Modified > I-UM-S?
    s = "v3h12: Last-Modified > I-UM-S?"
  case v3i4: // Moved permanently?
    s = "v3i4: Moved permanently?"
  case v3i7: // PUT?
    s = "v3i7: PUT?"
  case v3i12: // If-none-match exists?
    s = "v3i12: If-none-match exists?"
  case v3i13: // If-None-Match: * exists?
    s = "v3i13: If-None-Match: * exists?"
  case v3j18: // GET or HEAD?
    s = "v3j18: GET or HEAD?"
  case v3k5: // Moved permanently?
    s = "v3k5: Moved permanently?"
  case v3k7: // Previously existed?
    s = "v3k7: Previously existed?"
  case v3k13: // Etag in if-none-match?
    s = "v3k13: Etag in if-none-match?"
  case v3l5: // Moved temporarily?
    s = "v3l5: Moved temporarily?"
  case v3l7: // POST?
    s = "v3l7: POST?"
  case v3l13: // IMS exists?
    s = "v3l13: IMS exists?"
  case v3l14: // IMS is valid date?
    s = "v3l14: IMS is valid date?"
  case v3l15: // IMS > Now?
    s = "v3l15: IMS > Now?"
  case v3l17: // Last-Modified > IMS?
    s = "v3l17: Last-Modified > IMS?"
  case v3m5: // POST?
    s = "v3m5: POST?"
  case v3m7: // Server allows POST to missing resource?
    s = "v3m7: Server allows POST to missing resource?"
  case v3m16: // DELETE?
    s = "v3m16: DELETE?"
  case v3m20: // DELETE enacted immediately?
    s = "v3m20: DELETE enacted immediately?"
  case v3m20b: //
    s = "v3m20b: DELETE enacted immediately?"
  case v3n5: // Server allows POST to missing resource?
    s = "v3n5: Server allows POST to missing resource?"
  case v3n11: // Redirect?
    s = "v3n11: Redirect?"
  case v3n16: // POST?
    s = "v3n16: POST?"
  case v3o14: // Conflict?
    s = "v3o14: Conflict?"
  case v3o16: // PUT?
    s = "v3o16: PUT?"
  case v3o18: // Multiple representations?
    s = "v3o18: Multiple representations?"
  case v3o18b: //
    s = "v3o18b: Multiple representations?"
  case v3o20: // Response includes an entity?
    s = "v3o20: Response includes an entity?"
  case v3p3: // Conflict?
    s = "v3p3: Conflict?"
  case v3p11: // New resource?
    s = "v3p11: New resource?"
  default:
    s = "unknown decision"
  }
  return s
}


