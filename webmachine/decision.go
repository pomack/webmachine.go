package webmachine

import (
    "bytes"
    "compress/flate"
    "compress/gzip"
    "compress/lzw"
    "compress/zlib"
    "container/vector"
    "http"
    "io"
    "log"
    "mime"
    "os"
    "strconv"
    "strings"
    "time"
)

type wmDecisionCore struct {
    req                    Request
    resp                   ResponseWriter
    cxt                    Context
    handler                RequestHandler
    currentDecisionId      WMDecision
    lastModified           *time.Time
    unmodifiedSince        *time.Time
    modifiedSince          *time.Time
    mediaTypeOutputHandler MediaTypeHandler
    encodingOutputHandler  EncodingHandler
    charsetOutputHandler   CharsetHandler
    charsetInputHandler    CharsetHandler
    mediaType              string
    encoding               string
    charset                string
    language               string
    decisions              vector.IntVector
}

func handleRequest(handler RequestHandler, req Request, resp ResponseWriter) {
    d := &wmDecisionCore{req: req, resp: resp, handler: handler, currentDecisionId: v3b13}
    log.Print("[WM] Handling request for: ", req.Method(), " ", req.URL().Path, "\n")
    defer func() {
        log.Print("[WM] Running deferred function for: ", req.Method(), " ", req.URL().Path, "\n")
        go handler.FinishRequest(d.req, d.cxt)
        /*
           if e := recover(); e != nil {
             resp.WriteHeader(http.StatusInternalServerError)
           }
        */
    }()
    d.req, d.cxt = handler.StartRequest(d.req, d.cxt)
    nextDecision := v3b13
    log.Print("[WM] decision: ", nextDecision, " for ", req.Method(), " ", req.URL().Path, "\n")
    for nextDecision != wmResponded {
        nextDecision = d.makeDecision(nextDecision)
        log.Print("[WM] nextDecision: ", nextDecision, " for ", req.Method(), " ", req.URL().Path, "\n")
    }
}

func (p *wmDecisionCore) makeDecision(decisionId WMDecision) WMDecision {
    log.Print("[WM] Running decision: ", decisionId, " for ", p.req.Method(), " ", p.req.URL().Path, "\n")
    if decisionId != wmResponded {
        p.decisions.Push(int(decisionId))
    }
    p.currentDecisionId = decisionId
    p.logDecision(decisionId)
    nextDecision := p.decision(decisionId)
    if nextDecision != wmResponded {
        p.currentDecisionId = nextDecision
    }
    return nextDecision
}

func (p *wmDecisionCore) logDecision(decisionId WMDecision) {
    // TODO add logging
}

func (p *wmDecisionCore) writeHaltOrError(httpCode int, httpError os.Error) {
    p.resp.WriteHeader(httpCode)
    if httpError != nil {
        io.WriteString(p.resp, httpError.String())
    }
}

func (p *wmDecisionCore) decision(decisionId WMDecision) WMDecision {
    var nextDecision WMDecision
    switch decisionId {
    case v3b13:
        nextDecision = p.doV3b13() // Service available?
    case v3b13b:
        nextDecision = p.doV3b13b()
    case v3b12:
        nextDecision = p.doV3b12() // Known method?
    case v3b11:
        nextDecision = p.doV3b11() // URI too long?
    case v3b10:
        nextDecision = p.doV3b10() // Method allowed?
    case v3b9:
        nextDecision = p.doV3b9() // Malformed?
    case v3b8:
        nextDecision = p.doV3b8() // Authorized?
    case v3b7:
        nextDecision = p.doV3b7() // Forbidden?
    case v3b6:
        nextDecision = p.doV3b6() // Okay Content-* Headers?
    case v3b5:
        nextDecision = p.doV3b5() // Known Content-Type?
    case v3b4:
        nextDecision = p.doV3b4() // Req Entity Too Large?
    case v3b3:
        nextDecision = p.doV3b3() // OPTIONS?
    case v3c3:
        nextDecision = p.doV3c3() // Accept exists?
    case v3c4:
        nextDecision = p.doV3c4() // Acceptable media type available?
    case v3d4:
        nextDecision = p.doV3d4() // Accept-Language exists?
    case v3d5:
        nextDecision = p.doV3d5() // Acceptable Language available?
    case v3e5:
        nextDecision = p.doV3e5() // Accept-Charset exists?
    case v3e6:
        nextDecision = p.doV3e6() // Accceptable Charset available?
    case v3f6:
        nextDecision = p.doV3f6() // Accept-Encoding exists?
    case v3f7:
        nextDecision = p.doV3f7() // Acceptable encoding available?
    case v3g7:
        nextDecision = p.doV3g7() // Resource exists?
    case v3g8:
        nextDecision = p.doV3g8() // If-Match exists?
    case v3g9:
        nextDecision = p.doV3g9() // If-Match: * exists
    case v3g11:
        nextDecision = p.doV3g11() // ETag in If-Match
    case v3h7:
        nextDecision = p.doV3h7() // If-Match: * exists
    case v3h10:
        nextDecision = p.doV3h10() // If-unmodified-since exists?
    case v3h11:
        nextDecision = p.doV3h11() // I-UM-S is valid date?
    case v3h12:
        nextDecision = p.doV3h12() // Last-Modified > I-UM-S?
    case v3i4:
        nextDecision = p.doV3i4() // Moved permanently?
    case v3i7:
        nextDecision = p.doV3i7() // PUT?
    case v3i12:
        nextDecision = p.doV3i12() // If-none-match exists?
    case v3i13:
        nextDecision = p.doV3i13() // If-None-Match: * exists?
    case v3j18:
        nextDecision = p.doV3j18() // GET or HEAD?
    case v3k5:
        nextDecision = p.doV3k5() // Moved permanently?
    case v3k7:
        nextDecision = p.doV3k7() // Previously existed?
    case v3k13:
        nextDecision = p.doV3k13() // Etag in if-none-match?
    case v3l5:
        nextDecision = p.doV3l5() // Moved temporarily?
    case v3l7:
        nextDecision = p.doV3l7() // POST?
    case v3l13:
        nextDecision = p.doV3l13() // IMS exists?
    case v3l14:
        nextDecision = p.doV3l14() // IMS is valid date?
    case v3l15:
        nextDecision = p.doV3l15() // IMS > Now?
    case v3l17:
        nextDecision = p.doV3l17() // Last-Modified > IMS?
    case v3m5:
        nextDecision = p.doV3m5() // POST?
    case v3m7:
        nextDecision = p.doV3m7() // Server allows POST to missing resource?
    case v3m16:
        nextDecision = p.doV3m16() // DELETE?
    case v3m20:
        nextDecision = p.doV3m20() // DELETE enacted immediately?
    case v3m20b:
        nextDecision = p.doV3m20b() //
    case v3n5:
        nextDecision = p.doV3n5() // Server allows POST to missing resource?
    case v3n11:
        nextDecision = p.doV3n11() // Redirect?
    case v3n16:
        nextDecision = p.doV3n16() // POST?
    case v3o14:
        nextDecision = p.doV3o14() // Conflict?
    case v3o16:
        nextDecision = p.doV3o16() // PUT?
    case v3o18:
        nextDecision = p.doV3o18() // Multiple representations?
    case v3o20:
        nextDecision = p.doV3o20() // Response includes an entity?
    case v3p3:
        nextDecision = p.doV3p3() // Conflict?
    case v3p11:
        nextDecision = p.doV3p11() // New resource?
    default:
        p.resp.WriteHeader(501)
        nextDecision = wmResponded
    }
    return nextDecision
}

// Service Available
func (p *wmDecisionCore) doV3b13() WMDecision {
    if p.req == nil || p.handler == nil || p.resp == nil {
        if p.resp != nil {
            p.resp.WriteHeader(503)
        }
        return wmResponded
    }
    return v3b13b
}

func (p *wmDecisionCore) doV3b13b() WMDecision {
    var available bool
    var httpCode int
    var httpError os.Error
    if available, p.req, p.cxt, httpCode, httpError = p.handler.ServiceAvailable(p.req, p.cxt); available {
        return v3b12
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(503)
    return wmResponded
}

// Known method?
func (p *wmDecisionCore) doV3b12() WMDecision {
    method := p.req.Method()
    for _, m := range ALL_METHODS {
        if m == method {
            return v3b11
        }
    }
    p.resp.WriteHeader(501)
    return wmResponded
}

// URI too long?
func (p *wmDecisionCore) doV3b11() WMDecision {
    var tooLong bool
    var httpCode int
    var httpError os.Error
    if tooLong, p.req, p.cxt, httpCode, httpError = p.handler.URITooLong(p.req, p.cxt); tooLong {
        p.resp.WriteHeader(414)
        return wmResponded
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3b10
}

// Method allowed?
func (p *wmDecisionCore) doV3b10() WMDecision {
    var allowedMethods []string
    var httpCode int
    var httpError os.Error
    method := p.req.Method()
    allowedMethods, p.req, p.cxt, httpCode, httpError = p.handler.AllowedMethods(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    for _, allowedMethod := range allowedMethods {
        if method == allowedMethod {
            return v3b9
        }
    }
    p.resp.WriteHeader(405)
    s := "ALLOW " + strings.Join(allowedMethods, ",")
    // TODO handle error
    p.resp.Write([]byte(s))
    return wmResponded
}

// Malformed?
func (p *wmDecisionCore) doV3b9() WMDecision {
    var isMalformed bool
    var httpCode int
    var httpError os.Error
    if isMalformed, p.req, p.cxt, httpCode, httpError = p.handler.MalformedRequest(p.req, p.cxt); isMalformed {
        p.resp.WriteHeader(400)
        return wmResponded
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3b8
}

// Authorized?
func (p *wmDecisionCore) doV3b8() WMDecision {
    var isAuthorized bool
    var authHeaderString string
    var httpCode int
    var httpError os.Error
    if isAuthorized, authHeaderString, p.req, p.cxt, httpCode, httpError = p.handler.IsAuthorized(p.req, p.cxt); isAuthorized {
        return v3b7
    } else if len(authHeaderString) > 0 {
        p.resp.Header().Set("WWW-Authenticate", authHeaderString)
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(401)
    return wmResponded
}

// Forbidden?
func (p *wmDecisionCore) doV3b7() WMDecision {
    var forbidden bool
    var httpCode int
    var httpError os.Error
    if forbidden, p.req, p.cxt, httpCode, httpError = p.handler.Forbidden(p.req, p.cxt); forbidden {
        p.resp.WriteHeader(403)
        return wmResponded
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3b6
}

// Okay Content-* Headers?
func (p *wmDecisionCore) doV3b6() WMDecision {
    var isValid bool
    var httpCode int
    var httpError os.Error
    if isValid, p.req, p.cxt, httpCode, httpError = p.handler.ValidContentHeaders(p.req, p.cxt); isValid {
        return v3b5
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(501)
    return wmResponded
}

// Known Content-Type?
func (p *wmDecisionCore) doV3b5() WMDecision {
    var isKnown bool
    var httpCode int
    var httpError os.Error
    if isKnown, p.req, p.cxt, httpCode, httpError = p.handler.KnownContentType(p.req, p.cxt); isKnown {
        return v3b4
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(415)
    return wmResponded
}

// Req Entity Too Large?
func (p *wmDecisionCore) doV3b4() WMDecision {
    var isValid bool
    var httpCode int
    var httpError os.Error
    if isValid, p.req, p.cxt, httpCode, httpError = p.handler.ValidEntityLength(p.req, p.cxt); isValid {
        return v3b3
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(413)
    return wmResponded
}

// OPTIONS?
func (p *wmDecisionCore) doV3b3() WMDecision {
    var arr []string
    var httpCode int
    var httpError os.Error
    if p.req.Method() == OPTIONS {
        arr, p.req, p.cxt, httpCode, httpError = p.handler.Options(p.req, p.cxt)
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
        p.resp.WriteHeader(200)
        s := strings.Join(arr, ",")
        // TODO handle write error
        io.WriteString(p.resp, s)
        return wmResponded
    }
    return v3c3
}

// Accept exists?
func (p *wmDecisionCore) doV3c3() WMDecision {
    var provided []MediaTypeHandler
    var httpCode int
    var httpError os.Error
    arr, ok := p.req.Header()["Accept"]
    if !ok || len(arr) <= 0 {
        provided, p.req, p.cxt, httpCode, httpError = p.handler.ContentTypesProvided(p.req, p.cxt)
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
        if len(provided) >= 1 {
            p.mediaTypeOutputHandler = provided[0]
            p.resp.Header().Set("TCN", "choice")
            p.resp.Header().Set("Vary", "negotiate,accept")
        } else {
            // TODO Default is "text/html" and to_html
            p.mediaTypeOutputHandler = provided[0]
        }
        p.mediaType = p.mediaTypeOutputHandler.MediaType()
        p.resp.Header().Set("Content-Type", p.mediaType)
        _, params := mime.ParseMediaType(p.mediaType)
        if charset, ok := params["charset"]; ok {
            p.charset = charset
        }
        return v3d4
    }
    return v3c4
}

// Acceptable media type available?
func (p *wmDecisionCore) doV3c4() WMDecision {
    var provided []MediaTypeHandler
    var httpCode int
    var httpError os.Error
    arr, _ := p.req.Header()["Accept"]
    provided, p.req, p.cxt, httpCode, httpError = p.handler.ContentTypesProvided(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    mediaTypesProvided := make([]string, len(provided))
    for i, mth := range provided {
        mediaTypesProvided[i] = mth.MediaType()
    }
    bestMatch := chooseMediaType(mediaTypesProvided, arr[0])
    log.Print("[WDC]: Chose Media Type \"", bestMatch, "\" with accept ", arr, " and provided ", mediaTypesProvided)
    if len(bestMatch) > 0 {
        mediaType := bestMatch
        p.resp.Header().Set("Content-Type", mediaType)
        p.mediaType = mediaType
        _, params := mime.ParseMediaType(mediaType)
        if charset, ok := params["charset"]; ok {
            p.charset = charset
        }
        for _, mth := range provided {
            if mediaType == mth.MediaType() {
                p.mediaTypeOutputHandler = mth
                break
            }
        }
        return v3d4
    }
    p.resp.WriteHeader(406)
    return wmResponded
}

// Accept-Language exists?
func (p *wmDecisionCore) doV3d4() WMDecision {
    if arr, ok := p.req.Header()["Accept-Language"]; ok && len(arr) > 0 {
        return v3d5
    }
    return v3e5
}

// Acceptable Language available?
func (p *wmDecisionCore) doV3d5() WMDecision {
    var hasLanguage bool
    var httpCode int
    var httpError os.Error
    arr, _ := p.req.Header()["Accept-Language"]
    hasLanguage, p.req, p.cxt, httpCode, httpError = p.handler.IsLanguageAvailable(arr, p.req, p.cxt)
    if hasLanguage {
        p.language = arr[0]
        return v3e5
    } else if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(406)
    return wmResponded
}

// Accept-Charset exists?
func (p *wmDecisionCore) doV3e5() WMDecision {
    if arr, ok := p.req.Header()["Accept-Charset"]; ok && len(arr) > 0 {
        return v3e6
    }
    var handlers []CharsetHandler
    var httpCode int
    var httpError os.Error
    arr := make([]string, 1)
    arr[0] = "*"
    handlers, p.req, p.cxt, httpCode, httpError = p.handler.CharsetsProvided(arr, p.req, p.cxt)
    log.Print("Charsets Provided: ", handlers, " ", httpCode, " ", httpError, "\n")
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if len(handlers) == 0 {
        p.resp.WriteHeader(406)
        return wmResponded
    }
    p.charset = handlers[0].Charset()
    p.charsetOutputHandler = handlers[0]
    return v3f6
}

// Acceptable Charset available?
func (p *wmDecisionCore) doV3e6() WMDecision {
    arr, _ := p.req.Header()["Accept-Charset"]
    var handlers []CharsetHandler
    var httpCode int
    var httpError os.Error
    handlers, p.req, p.cxt, httpCode, httpError = p.handler.CharsetsProvided(arr, p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if len(handlers) == 0 {
        p.resp.WriteHeader(406)
        return wmResponded
    }
    p.charset = handlers[0].Charset()
    p.charsetOutputHandler = handlers[0]
    return v3f6
}

// Accept-Encoding exists?
func (p *wmDecisionCore) doV3f6() WMDecision {
    ctype := p.mediaTypeOutputHandler
    cset := p.charsetOutputHandler
    cs := ""
    if cset != nil {
        cs = cset.Charset()
        if len(cs) > 0 {
            cs = "; charset=" + cs
        }
    }
    headers := p.resp.Header()
    headers.Set("Content-Type", ctype.MediaType()+cs)
    if arr, ok := p.req.Header()["Accept-Encoding"]; ok && len(arr) > 0 {
        return v3f7
    }
    arr := make([]string, 1)
    arr[0] = "identity;q=1.0,*;q=0.5"
    var handlers []EncodingHandler
    var httpCode int
    var httpError os.Error
    handlers, p.req, p.cxt, httpCode, httpError = p.handler.EncodingsProvided(arr, p.req, p.cxt)
    if len(handlers) > 0 {
        p.encodingOutputHandler = handlers[0]
        p.encoding = handlers[0].Encoding()
        return v3g7
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    p.resp.WriteHeader(406)
    return wmResponded
    return v3g7
}

// Acceptable Encoding available?
func (p *wmDecisionCore) doV3f7() WMDecision {
    if len(p.chooseEncoding()) == 0 {
        p.resp.WriteHeader(406)
        return wmResponded
    }
    return v3g7
    /*
       arr, _ := p.req.Header()["Accept-Encoding"]
       var handlers []EncodingHandler
       var httpCode int
       var httpError os.Error
       handlers, p.req, p.cxt, httpCode, httpError = p.handler.EncodingsProvided(arr, p.req, p.cxt)
       log.Print("[WDC]: Accept-Encoding: \"", arr, "\" vs handlers: ", handlers)
       if len(handlers) > 0 {
         p.encodingOutputHandler = handlers[0]
         p.encoding = handlers[0].Encoding()
         return v3g7
       }
       if httpCode > 0 {
         p.writeHaltOrError(httpCode, httpError)
         return wmResponded
       }
       p.resp.WriteHeader(406)
       return wmResponded
    */
}

// Resource exists?
func (p *wmDecisionCore) doV3g7() WMDecision {
    variances := p.variances()
    if len(variances) > 0 {
        p.resp.Header().Set("Vary", strings.Join(variances, ", "))
    }
    var exists bool
    var httpCode int
    var httpError os.Error
    exists, p.req, p.cxt, httpCode, httpError = p.handler.ResourceExists(p.req, p.cxt)
    if exists {
        return v3g8
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3h7
}

// If-Match exists?
func (p *wmDecisionCore) doV3g8() WMDecision {
    if arr, ok := p.req.Header()["If-Match"]; ok && len(arr) > 0 {
        return v3g9
    }
    return v3h10
}

// If-Match: * exists
func (p *wmDecisionCore) doV3g9() WMDecision {
    if arr, ok := p.req.Header()["If-Match"]; ok && len(arr) > 0 && arr[0] == "*" {
        return v3h10
    }
    return v3g11
}

// ETag in If-Match
func (p *wmDecisionCore) doV3g11() WMDecision {
    if arr, ok := p.req.Header()["If-Match"]; ok && len(arr) > 0 {
        etag, err := strconv.Unquote(arr[0])
        if err != nil {
            etag = arr[0]
        }
        var generatedEtag string
        var httpCode int
        var httpError os.Error
        generatedEtag, p.req, p.cxt, httpCode, httpError = p.handler.GenerateETag(p.req, p.cxt)
        if generatedEtag == etag {
            return v3h10
        }
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
    }
    p.resp.WriteHeader(412)
    return wmResponded
}

// If-Match: * exists
func (p *wmDecisionCore) doV3h7() WMDecision {
    if arr, ok := p.req.Header()["If-Match"]; ok && len(arr) > 0 && arr[0] == "*" {
        p.resp.WriteHeader(412)
        return wmResponded
    }
    return v3i7
}

// If-unmodified-since exists?
func (p *wmDecisionCore) doV3h10() WMDecision {
    if arr, ok := p.req.Header()["If-Unmodified-Since"]; ok && len(arr) > 0 {
        return v3h11
    }
    return v3i12
}

// I-UM-S is valid date?
func (p *wmDecisionCore) doV3h11() WMDecision {
    arr, _ := p.req.Header()["If-Unmodified-Since"]
    iumsDate := arr[0]
    t, err := time.Parse(http.TimeFormat, iumsDate)
    p.unmodifiedSince = t
    if err == nil {
        return v3h12
    }
    return v3i12
}

// Last-Modified > I-UM-S?
func (p *wmDecisionCore) doV3h12() WMDecision {
    var lastModified *time.Time
    var httpCode int
    var httpError os.Error
    lastModified, p.req, p.cxt, httpCode, httpError = p.handler.LastModified(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    t := p.unmodifiedSince
    if t != nil && lastModified != nil {
        log.Print("[WMC]: comparing Last-Modified internal: ", t.Seconds(), " vs. received from client ", lastModified.Seconds())
    }
    if t != lastModified && t != nil && lastModified != nil && lastModified.Seconds() > t.Seconds() {
        p.resp.WriteHeader(412)
        return wmResponded
    }
    return v3i12
}

// Moved permanently? (apply PUT to different URI)
func (p *wmDecisionCore) doV3i4() WMDecision {
    var uri string
    var httpCode int
    var httpError os.Error
    uri, p.req, p.cxt, httpCode, httpError = p.handler.MovedPermanently(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if len(uri) > 0 {
        p.resp.Header().Set("Location", uri)
        p.resp.WriteHeader(301)
        return wmResponded
    }
    return v3p3
}

// PUT?
func (p *wmDecisionCore) doV3i7() WMDecision {
    if p.req.Method() == PUT {
        return v3i4
    }
    return v3k7
}

// If-None-Match exists?
func (p *wmDecisionCore) doV3i12() WMDecision {
    if arr, ok := p.req.Header()["If-None-Match"]; ok && len(arr) > 0 {
        return v3i13
    }
    return v3l13
}

// If-None-Match: * exists
func (p *wmDecisionCore) doV3i13() WMDecision {
    if arr, ok := p.req.Header()["If-None-Match"]; ok && len(arr) > 0 && arr[0] == "*" {
        return v3j18
    }
    return v3k13
}

// GET or HEAD?
func (p *wmDecisionCore) doV3j18() WMDecision {
    method := p.req.Method()
    if method == GET || method == HEAD {
        p.resp.WriteHeader(304)
        return wmResponded
    }
    p.resp.WriteHeader(412)
    return wmResponded
}

// Moved permanently?
func (p *wmDecisionCore) doV3k5() WMDecision {
    var uri string
    var httpCode int
    var httpError os.Error
    uri, p.req, p.cxt, httpCode, httpError = p.handler.MovedPermanently(p.req, p.cxt)
    if len(uri) > 0 {
        p.resp.Header().Set("Location", uri)
        p.resp.WriteHeader(301)
        return wmResponded
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3l5
}

// Previously existed?
func (p *wmDecisionCore) doV3k7() WMDecision {
    var previouslyExisted bool
    var httpCode int
    var httpError os.Error
    previouslyExisted, p.req, p.cxt, httpCode, httpError = p.handler.PreviouslyExisted(p.req, p.cxt)
    if previouslyExisted {
        return v3k5
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3l7
}

// ETag in If-None-Match
func (p *wmDecisionCore) doV3k13() WMDecision {
    if arr, ok := p.req.Header()["If-None-Match"]; ok && len(arr) > 0 {
        etag, err := strconv.Unquote(arr[0])
        if err != nil {
            etag = arr[0]
        }
        var generatedEtag string
        var httpCode int
        var httpError os.Error
        generatedEtag, p.req, p.cxt, httpCode, httpError = p.handler.GenerateETag(p.req, p.cxt)
        if generatedEtag == etag {
            return v3j18
        }
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
    }
    return v3l13
}

// Moved permanently? (apply PUT to different URI)
func (p *wmDecisionCore) doV3l5() WMDecision {
    var uri string
    var httpCode int
    var httpError os.Error
    uri, p.req, p.cxt, httpCode, httpError = p.handler.MovedTemporarily(p.req, p.cxt)
    if len(uri) > 0 {
        p.resp.Header().Set("Location", uri)
        p.resp.WriteHeader(307)
        return wmResponded
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    return v3m5
}

// POST?
func (p *wmDecisionCore) doV3l7() WMDecision {
    if p.req.Method() == POST {
        return v3m7
    }
    log.Print("Writing Header 404\n")
    p.resp.WriteHeader(404)
    log.Print("Done Writing Header 404\n")
    return wmResponded
}

// If-Modified-Since exists?
func (p *wmDecisionCore) doV3l13() WMDecision {
    if arr, ok := p.req.Header()["If-Modified-Since"]; ok && len(arr) > 0 {
        return v3l14
    }
    return v3m16
}

// I-M-S is valid date?
func (p *wmDecisionCore) doV3l14() WMDecision {
    arr, _ := p.req.Header()["If-Modified-Since"]
    iumsDate := arr[0]
    t, err := time.Parse(http.TimeFormat, iumsDate)
    p.modifiedSince = t
    if err == nil && t != nil {
        return v3l15
    }
    return v3m16
}

// I-M-S > Now?
func (p *wmDecisionCore) doV3l15() WMDecision {
    now := time.UTC().Seconds()
    t := p.modifiedSince
    if t != nil && t.Seconds() > now {
        return v3m16
    }
    return v3l17
}

// Last-Modified > I-M-S?
func (p *wmDecisionCore) doV3l17() WMDecision {
    t := p.modifiedSince
    var lastModified *time.Time
    var httpCode int
    var httpError os.Error
    lastModified, p.req, p.cxt, httpCode, httpError = p.handler.LastModified(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if lastModified == nil || t == nil || lastModified.Seconds() > t.Seconds() {
        return v3m16
    }
    p.resp.WriteHeader(304)
    return wmResponded
}

// POST?
func (p *wmDecisionCore) doV3m5() WMDecision {
    if p.req.Method() == POST {
        return v3n5
    }
    p.resp.WriteHeader(410)
    return wmResponded
}

// Server allows POST to missing resource?
func (p *wmDecisionCore) doV3m7() WMDecision {
    var allowMissingPost bool
    var httpCode int
    var httpError os.Error
    allowMissingPost, p.req, p.cxt, httpCode, httpError = p.handler.AllowMissingPost(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if allowMissingPost {
        return v3n11
    }
    p.resp.WriteHeader(404)
    return wmResponded
}

// DELETE?
func (p *wmDecisionCore) doV3m16() WMDecision {
    if p.req.Method() == DELETE {
        return v3m20
    }
    return v3n16
}

// DELETE enacted immediately?
// Also where DELETE is forced
func (p *wmDecisionCore) doV3m20() WMDecision {
    var ok bool
    var httpCode int
    var httpError os.Error
    ok, p.req, p.cxt, httpCode, httpError = p.handler.DeleteResource(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if ok {
        return v3m20b
    }
    p.resp.WriteHeader(500)
    return wmResponded
}

// Check if totally deleted?
func (p *wmDecisionCore) doV3m20b() WMDecision {
    var completed bool
    var httpCode int
    var httpError os.Error
    completed, p.req, p.cxt, httpCode, httpError = p.handler.DeleteCompleted(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if completed {
        return v3o20
    }
    p.resp.WriteHeader(202)
    return wmResponded
}

// Server allows POST to missing resource?
func (p *wmDecisionCore) doV3n5() WMDecision {
    var allowed bool
    var httpCode int
    var httpError os.Error
    allowed, p.req, p.cxt, httpCode, httpError = p.handler.AllowMissingPost(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if allowed {
        return v3n11
    }
    p.resp.WriteHeader(410)
    return wmResponded
}

// Redirect -- only accessible if method == POST
func (p *wmDecisionCore) doV3n11() WMDecision {
    var postIsCreate bool
    var httpCode int
    var httpError os.Error
    postIsCreate, p.req, p.cxt, httpCode, httpError = p.handler.PostIsCreate(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if postIsCreate {
        log.Print("[WM]: v3n11: Creating Path\n")
        _, p.req, p.cxt, httpCode, httpError = p.handler.CreatePath(p.req, p.cxt)
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
        log.Print("[WM]: v3n11: Running Accept Helper\n")
        n, err := p.acceptHelper()
        if err != nil {
            p.resp.WriteHeader(n)
            io.WriteString(p.resp, err.String())
            return wmResponded
        }
        log.Print("[WM]: v3n11: Running Process Post\n")
        _, p.req, p.cxt, httpCode, httpError = p.handler.ProcessPost(p.req, p.cxt)
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
        // TODO Aalok check what should be done here
        //p.resp.WriteHeader(n)
        //log.Print("Wrote Header but may not return wmResponded in doV3n11()\n")
        //p.encodeBodyIfSet()
    }
    log.Print("[WM]: v3n11: Running Response is Redirect?\n")
    var respIsRedirect bool
    respIsRedirect, p.req, p.cxt, httpCode, httpError = p.handler.ResponseIsRedirect(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if respIsRedirect {
        if len(p.resp.Header().Get("Location")) > 0 {
            p.resp.WriteHeader(303)
        } else {
            p.resp.WriteHeader(500)
            io.WriteString(p.resp, "Response had do_redirect but no Location")
        }
        return wmResponded
    }
    return v3p11
}

// POST?
func (p *wmDecisionCore) doV3n16() WMDecision {
    if p.req.Method() == POST {
        return v3n11
    }
    return v3o16
}

// Conflict?
func (p *wmDecisionCore) doV3o14() WMDecision {
    var isConflict bool
    var httpCode int
    var httpError os.Error
    // TOOD v3n11
    isConflict, p.req, p.cxt, httpCode, httpError = p.handler.IsConflict(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if isConflict {
        p.resp.WriteHeader(409)
        return wmResponded
    }
    code, err := p.acceptHelper()
    if err != nil {
        p.resp.WriteHeader(code)
        io.WriteString(p.resp, err.String())
        return wmResponded
    }
    if code > 0 {
        p.resp.WriteHeader(code)
        return wmResponded
    }
    return v3p11
}

// PUT?
func (p *wmDecisionCore) doV3o16() WMDecision {
    if p.req.Method() == PUT {
        return v3o14
    }
    return v3o18
}

// Multiple representations?
// Also where body generation for GET and HEAD is done
func (p *wmDecisionCore) doV3o18() WMDecision {
    method := p.req.Method()
    buildBody := method == GET || method == HEAD
    var multipleChoices bool
    var httpCode int
    var httpError os.Error
    var httpHeaders http.Header
    multipleChoices, httpHeaders, p.req, p.cxt, httpCode, httpError = p.handler.MultipleChoices(p.req, p.cxt)
    if httpHeaders != nil {
        headers := p.resp.Header()
        for k, v := range httpHeaders {
            if headers.Get(k) != "" {
                for _, v1 := range v {
                    headers.Set(k, v1)
                }
            } else {
                for _, v1 := range v {
                    headers.Add(k, v1)
                }
            }
        }
    }
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if buildBody {
        var etag string
        var httpCode int
        var httpError os.Error
        etag, p.req, p.cxt, httpCode, httpError = p.handler.GenerateETag(p.req, p.cxt)
        if httpCode > 0 {
            p.writeHaltOrError(httpCode, httpError)
            return wmResponded
        }
        if len(etag) > 0 {
            p.resp.Header().Set("ETag", strconv.Quote(etag))
        }
        var lastModified, expires *time.Time
        lastModified, p.req, p.cxt, httpCode, httpError = p.handler.LastModified(p.req, p.cxt)
        if lastModified != nil {
            p.resp.Header().Set("Last-Modified", lastModified.Format(http.TimeFormat))
        }
        expires, p.req, p.cxt, httpCode, httpError = p.handler.Expires(p.req, p.cxt)
        if expires != nil {
            p.resp.Header().Set("Expires", expires.Format(http.TimeFormat))
        }
        if p.mediaTypeOutputHandler != nil {
            p.mediaTypeOutputHandler.OutputTo(p.req, p.cxt, p.resp, p.resp)
            p.resp.Flush()
            return wmResponded
        } else {
            var provided []MediaTypeHandler
            provided, p.req, p.cxt, httpCode, httpError = p.handler.ContentTypesProvided(p.req, p.cxt)
            if httpCode > 0 {
                p.writeHaltOrError(httpCode, httpError)
                return wmResponded
            }
            if len(provided) > 0 && len(provided) == 1 {
                provided[0].OutputTo(p.req, p.cxt, p.resp, p.resp)
                return wmResponded
            }
        }
    }
    if multipleChoices {
        p.resp.WriteHeader(300)
        return wmResponded
    }
    p.resp.WriteHeader(200)
    return wmResponded
}

// Redirect
func (p *wmDecisionCore) doV3o20() WMDecision {
    if p.handler.HasRespBody(p.req, p.cxt) {
        return v3o18
    }
    p.resp.WriteHeader(204)
    return wmResponded
}

// Conflict?
func (p *wmDecisionCore) doV3p3() WMDecision {
    var isConflict bool
    var httpCode int
    var httpError os.Error
    // TOOD v3n11
    isConflict, p.req, p.cxt, httpCode, httpError = p.handler.IsConflict(p.req, p.cxt)
    log.Print("[WDC]: V3P3: isConflict", isConflict, ", code: ", httpCode, ", error: ", httpError)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return wmResponded
    }
    if isConflict {
        p.resp.WriteHeader(409)
        return wmResponded
    }
    code, err := p.acceptHelper()
    log.Print("[WDC]: V3P3: acceptHelper code: ", code, ", err: ", err)
    if err != nil {
        p.resp.WriteHeader(code)
        io.WriteString(p.resp, err.String())
        return wmResponded
    }
    if code > 0 {
        p.resp.WriteHeader(code)
        return wmResponded
    }
    return v3p11
}

// New resource?  (at this point boils down to "has location header")
func (p *wmDecisionCore) doV3p11() WMDecision {
    if _, ok := p.resp.Header()["Location"]; ok {
        p.resp.WriteHeader(201)
        return wmResponded
    }
    return v3o20
}

func (p *wmDecisionCore) acceptHelper() (int, os.Error) {
    // TODO acceptHelper
    ct := p.req.Header().Get("Content-Type")
    if len(ct) == 0 {
        ct = MIME_TYPE_OCTET_STREAM
    }
    var ctAccepted []MediaTypeInputHandler
    var httpCode int
    var httpHeaders http.Header
    var httpError os.Error
    var buf *bytes.Buffer
    ctAccepted, p.req, p.cxt, httpCode, httpError = p.handler.ContentTypesAccepted(p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return httpCode, nil
    }
    arr := make([]string, len(ctAccepted))
    for i := 0; i < len(arr); i++ {
        arr[i] = ctAccepted[i].MediaType()
    }
    mt := chooseMediaType(arr, ct)
    acceptArr := splitAcceptString(mt)
    if len(acceptArr) == 0 {
        return 415, nil
    }
    for i := 0; i < len(arr); i++ {
        if arr[i] == mt {
            log.Print("[AH]: Capturing accepted value of type ", mt)
            buf = bytes.NewBuffer(make([]byte, 0))
            httpCode, httpHeaders, httpError = ctAccepted[i].OutputTo(p.req, p.cxt, buf)
            break
        }
    }
    if httpHeaders != nil {
        headers := p.resp.Header()
        for k, v := range httpHeaders {
            if headers.Get(k) != "" {
                for _, v1 := range v {
                    headers.Set(k, v1)
                }
            } else {
                for _, v1 := range v {
                    headers.Add(k, v1)
                }
            }
        }
    }
    log.Print("[AH]: Done capturing input stream of type ", mt)
    if httpError == nil {
        return httpCode, buf
    }
    return httpCode, httpError
}

func (p *wmDecisionCore) encodeBodyIfSet() bool {
    if !p.handler.HasRespBody(p.req, p.cxt) {
        return false
    }
    bodyWriter, _ := p.bodyEncoder(p.resp)
    p.mediaTypeOutputHandler.OutputTo(p.req, p.cxt, bodyWriter, p.resp)
    return true
}

func (p *wmDecisionCore) bodyEncoder(w io.Writer) (io.Writer, os.Error) {
    var outW io.Writer
    var err os.Error
    switch p.encoding {
    default:
        outW = w
    case "identity":
        outW = w
    case "deflate":
        outW = flate.NewWriter(w, 6)
    case "gzip":
        outW, err = gzip.NewWriter(w)
    case "lzw":
        outW = lzw.NewWriter(w, lzw.LSB, 8)
    case "zlib":
        outW, err = zlib.NewWriter(w)
    }
    return outW, err
}

func (p *wmDecisionCore) chooseEncoding() string {
    var encodingHandlers []EncodingHandler
    arr := make([]string, 1)
    arr[0] = "*"
    encodingHandlers, p.req, p.cxt, _, _ = p.handler.EncodingsProvided(arr, p.req, p.cxt)
    log.Print("[WDC]: chooseEncoding: ", encodingHandlers)
    if len(encodingHandlers) == 0 {
        return ""
    }
    encodingMap := make(map[string]EncodingHandler)
    encodings := make([]string, len(encodingHandlers))
    for i, encodingHandler := range encodingHandlers {
        encodings[i] = encodingHandler.Encoding()
        encodingMap[encodingHandler.Encoding()] = encodingHandler
    }
    s := p.req.Header().Get("Accept-Encoding")
    encoding := chooseEncodingWithDefaultString(encodings, s)
    log.Print("[WDC]: Accept-Encoding: \"", s, "\", we supply ", encodings, " and choosing encoding \"", encoding, "\"")
    if len(encoding) > 0 {
        p.resp.Header().Set("Content-Encoding", encoding)
        p.resp.AddEncoding(encodingMap[encoding], p.req, p.cxt)
        p.encodingOutputHandler = encodingMap[encoding]
    }
    return encoding
}

func (p *wmDecisionCore) chooseCharset() string {
    var charsetHandlers []CharsetHandler
    var httpCode int
    var httpError os.Error
    arr := make([]string, 1)
    arr[0] = "*"
    charsetHandlers, p.req, p.cxt, httpCode, httpError = p.handler.CharsetsProvided(arr, p.req, p.cxt)
    if httpCode > 0 {
        p.writeHaltOrError(httpCode, httpError)
        return ""
    }
    if len(charsetHandlers) <= 0 {
        return ""
    }
    charsetMap := make(map[string]CharsetHandler)
    charsets := make([]string, len(charsetHandlers))
    for i, charsetHandler := range charsetHandlers {
        charsets[i] = charsetHandler.Charset()
        charsetMap[charsetHandler.Charset()] = charsetHandler
    }
    s := p.req.Header().Get("Accept-Charset")
    charset := chooseCharsetWithDefaultString(charsets, s)
    if len(charset) > 0 {
        // TODO append to Content-Type header response if not already specified
        p.charset = charset
        p.charsetOutputHandler = charsetMap[charset]
    }
    return charset
}

func (p *wmDecisionCore) variances() []string {
    var v vector.StringVector
    var ctp []MediaTypeHandler
    var ep []EncodingHandler
    var cp []CharsetHandler
    arr := make([]string, 1)
    arr[0] = "*"
    ctp, p.req, p.cxt, _, _ = p.handler.ContentTypesProvided(p.req, p.cxt)
    ep, p.req, p.cxt, _, _ = p.handler.EncodingsProvided(arr, p.req, p.cxt)
    cp, p.req, p.cxt, _, _ = p.handler.CharsetsProvided(arr, p.req, p.cxt)
    if len(ctp) > 1 {
        v.Push("Accept")
    }
    if len(ep) > 1 {
        v.Push("Accept-Encoding")
    }
    if len(cp) > 1 {
        v.Push("Accept-Charset")
    }
    var headers []string
    headers, p.req, p.cxt, _, _ = p.handler.Variances(p.req, p.cxt)
    v2 := vector.StringVector(headers)
    v.AppendVector(&v2)
    return v
}
