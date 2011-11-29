package webmachine

import (
    "http"
    "io"
    "json"
    "log"
    "path"
    "os"
)


type jsonWriter struct {
    obj interface{}
}

func newJSONWriter(obj interface{}) *jsonWriter {
    return &jsonWriter{obj:obj}
}

func (p *jsonWriter) WriteTo(writer io.Writer) (n int64, err os.Error) {
    w := json.NewEncoder(writer)
    err = w.Encode(p.obj)
    return
}

func (p *jsonWriter) String() string {
    b, err := json.Marshal(p.obj)
    if err != nil {
        return err.String()
    }
    return string(b)
}


type PassThroughMediaTypeInputHandler struct {
    mediaType           string
    charset             string
    language            string
    filename            string
    urlPath             string
    append              bool
    numberOfBytes       int64
    reader              io.Reader
    writtenStatusHeader bool
}

func NewPassThroughMediaTypeInputHandler(mediaType, charset, language, filename, urlPath string, append bool, numberOfBytes int64, reader io.Reader) *PassThroughMediaTypeInputHandler {
    return &PassThroughMediaTypeInputHandler{
        mediaType:     mediaType,
        charset:       charset,
        language:      language,
        filename:      filename,
        urlPath:       urlPath,
        append:        append,
        numberOfBytes: numberOfBytes,
        reader:        reader,
    }
}

func (p *PassThroughMediaTypeInputHandler) MediaTypeInput() string {
    return p.mediaType
}

func (p *PassThroughMediaTypeInputHandler) MediaTypeHandleInputFrom(req Request, cxt Context) (int, http.Header, io.WriterTo) {
    fileInfo, err := os.Stat(p.filename)
    var file *os.File
    m := make(map[string]string)
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
            //headers.Set("Content-Type", MIME_TYPE_JSON)
            m["status"] = "error"
            m["message"] = err.String()
            m["result"] = p.urlPath
            return http.StatusInternalServerError, headers, newJSONWriter(m)
        }
        if file, err = os.OpenFile(p.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
            log.Print("[PTMTIH]: Unable to create file named: \"", p.filename, "\" due to error: ", err)
            headers := make(http.Header)
            //headers.Set("Content-Type", MIME_TYPE_JSON)
            m["status"] = "error"
            m["message"] = err.String()
            m["result"] = p.urlPath
            return http.StatusInternalServerError, headers, newJSONWriter(m)
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
            //headers.Set("Content-Type", MIME_TYPE_JSON)
            m["status"] = "error"
            m["message"] = err.String()
            m["result"] = p.urlPath
            return http.StatusInternalServerError, headers, newJSONWriter(m)
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
        //headers.Set("Content-Type", MIME_TYPE_JSON)
        m["status"] = "error"
        m["message"] = err.String()
        m["result"] = p.urlPath
        return http.StatusInternalServerError, headers, newJSONWriter(m)
    }
    headers := make(http.Header)
    //headers.Set("Content-Type", MIME_TYPE_JSON)
    m["status"] = "success"
    m["message"] = ""
    m["result"] = p.urlPath
    return http.StatusOK, headers, newJSONWriter(m)
}
