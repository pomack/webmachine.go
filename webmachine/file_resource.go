package webmachine

import (
    "container/list"
    "http"
    "io"
    "log"
    "path"
    "path/filepath"
    "os"
    "mime"
    "rand"
    "strconv"
    "time"
)

type FileResource struct {
    DefaultRequestHandler
    dirPath               string
    urlPathPrefix         string
    allowWrite            bool
    allowDirectoryListing bool
}

type FileResourceContext interface {
    FullPath() string
    SetFullPath(fullPath string)
    ReaderOpen() (io.ReadCloser, os.Error)
    WriterOpen(append bool) (io.WriteCloser, os.Error)
    Close() os.Error
    Read(data []byte) (int, os.Error)
    Write(data []byte) (int, os.Error)
    Exists() bool
    CanRead() bool
    CanWrite(append bool) bool
    IsDir() bool
    IsFile() bool
    LastModified() *time.Time
    Len() int64
    HasMultipleResources() bool
    MultipleResourceNames() []string
}

type fileResourceContext struct {
    fullPath       string
    fileInfo       *os.FileInfo
    namedResources []string
    reader         io.ReadCloser
    writer         io.WriteCloser
}

func NewFileResourceContext() FileResourceContext {
    return &fileResourceContext{namedResources: make([]string, 0)}
}

func NewFileResourceContextWithPath(fullPath string) FileResourceContext {
    cxt := NewFileResourceContext()
    cxt.SetFullPath(fullPath)
    return cxt
}

func (p *fileResourceContext) FullPath() string {
    return p.fullPath
}

func (p *fileResourceContext) SetFullPath(fullPath string) {
    p.fullPath = fullPath
    p.fileInfo, _ = os.Stat(fullPath)
    if len(p.namedResources) > 0 {
        p.namedResources = make([]string, 0)
    }
    if p.fileInfo == nil {
        dir, tail := path.Split(fullPath)
        dirInfo, _ := os.Stat(dir)
        if dirInfo != nil && dirInfo.IsDirectory() {
            dirFile, _ := os.Open(dir)
            if dirFile != nil {
                names, _ := dirFile.Readdirnames(-1)
                l := list.New()
                for _, name := range names {
                    if len(name) > len(tail) && name[len(tail)] == '.' && name[:len(tail)] == tail {
                        l.PushBack(name)
                    }
                }
                if l.Len() > 0 {
                    i := 0
                    namedResources := make([]string, l.Len())
                    for e := l.Front(); e != nil; e = e.Next() {
                        namedResources[i] = e.Value.(string)
                        i++
                    }
                    p.namedResources = namedResources
                }
            }
        }
    }
}

func (p *fileResourceContext) HasMultipleResources() bool {
    return len(p.namedResources) > 0
}

func (p *fileResourceContext) MultipleResourceNames() []string {
    return p.namedResources
}

func (p *fileResourceContext) ReaderOpen() (io.ReadCloser, os.Error) {
    if p.reader != nil {
        p.reader.Close()
    }
    var err os.Error
    p.reader, err = os.Open(p.fullPath)
    return p, err
}

func (p *fileResourceContext) WriterOpen(append bool) (io.WriteCloser, os.Error) {
    if p.writer != nil {
        p.writer.Close()
    }
    var err os.Error
    if append {
        p.writer, err = os.OpenFile(p.fullPath, os.O_APPEND, 0644)
    } else {
        p.writer, err = os.OpenFile(p.fullPath, os.O_WRONLY, 0644)
    }
    return p, err
}

func (p *fileResourceContext) Close() os.Error {
    var e1, e2 os.Error
    p.fileInfo = nil
    if p.reader != nil {
        e1 = p.reader.Close()
    }
    if p.writer != nil {
        e2 = p.writer.Close()
    }
    if e1 != nil {
        return e1
    }
    return e2
}

func (p *fileResourceContext) Read(data []byte) (int, os.Error) {
    if p.reader == nil {
        log.Print("[FRC]: Trying to open file ", p.FullPath(), " for reading\n")
        var err os.Error
        p.reader, err = os.Open(p.FullPath())
        if err != nil {
            return 0, err
        }
        if p.reader == nil {
            return 0, os.EOF
        }
    }
    log.Print("[FRC]: Going to read ", len(data), " bytes\n")
    return p.reader.Read(data)
}

func (p *fileResourceContext) Write(data []byte) (int, os.Error) {
    if p.writer == nil {
        log.Print("[FRC]: Trying to open file ", p.FullPath(), " for appending\n")
        var err os.Error
        p.writer, err = os.OpenFile(p.FullPath(), os.O_APPEND, 0644)
        if err != nil {
            return 0, err
        }
    }
    log.Print("[FRC]: Going to write ", len(data), " bytes\n")
    return p.writer.Write(data)
}

func (p *fileResourceContext) Exists() bool {
    return p.fileInfo != nil
}

func (p *fileResourceContext) CanRead() bool {
    if p.fileInfo == nil {
        return false
    }
    file, err := os.Open(p.fullPath)
    if file != nil {
        file.Close()
    }
    return err == nil
}

func (p *fileResourceContext) CanWrite(append bool) bool {
    file, err := os.OpenFile(p.fullPath, os.O_APPEND, 0644)
    if file != nil {
        file.Close()
    }
    return err == nil
}

func (p *fileResourceContext) IsDir() bool {
    return p.fileInfo != nil && p.fileInfo.IsDirectory()
}

func (p *fileResourceContext) IsFile() bool {
    return p.fileInfo != nil && p.fileInfo.IsRegular()
}

func (p *fileResourceContext) LastModified() *time.Time {
    if p.fileInfo != nil {
        return time.SecondsToUTC(int64(p.fileInfo.Mtime_ns / 1e9))
    }
    return nil
}

func (p *fileResourceContext) Len() int64 {
    if p.fileInfo != nil {
        return p.fileInfo.Size
    }
    return 0
}

func NewFileResource(directoryPath, urlPathPrefix string, allowWrite bool, allowDirectoryListing bool) *FileResource {
    return &FileResource{dirPath: directoryPath, urlPathPrefix: urlPathPrefix, allowWrite: allowWrite, allowDirectoryListing: allowDirectoryListing}
}

func (p *FileResource) GenerateContext(req Request, cxt Context) FileResourceContext {
    if frc, ok := cxt.(FileResourceContext); ok {
        return frc
    }
    fullPath := filepath.Join(p.dirPath, filepath.Clean(req.URL().Path[len(p.urlPathPrefix):]))
    return NewFileResourceContextWithPath(fullPath)
}

func (p *FileResource) HandlerFor(req Request, writer ResponseWriter) RequestHandler {
    path := req.URL().Path
    if path >= p.urlPathPrefix && path[0:len(p.urlPathPrefix)] == p.urlPathPrefix {
        return p
    }
    return nil
}

func (p *FileResource) StartRequest(req Request, cxt Context) (Request, Context) {
    frc := p.GenerateContext(req, cxt)
    frc.SetFullPath(filepath.Join(p.dirPath, filepath.Clean(req.URL().Path[len(p.urlPathPrefix):])))
    return req, frc
}
/*
func (p *FileResource) ServiceAvailable(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
  return true, req, cxt, 0, nil
}
*/
func (p *FileResource) ResourceExists(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    if !frc.Exists() && !frc.HasMultipleResources() {
        return false, req, frc, 0, nil
    }
    if frc.IsDir() {
        if req.Method() == GET || req.Method() == HEAD {
            return p.allowDirectoryListing, req, frc, 0, nil
        }
    }
    return frc.IsFile() || frc.HasMultipleResources(), req, frc, 0, nil
}

func (p *FileResource) AllowedMethods(req Request, cxt Context) ([]string, Request, Context, int, os.Error) {
    var methods []string
    if p.allowWrite {
        methods = []string{GET, HEAD, POST, PUT, DELETE}
    } else {
        methods = []string{GET, HEAD}
    }
    return methods, req, cxt, 0, nil
}

func (p *FileResource) IsAuthorized(req Request, cxt Context) (bool, string, Request, Context, int, os.Error) {
    method := req.Method()
    frc := cxt.(FileResourceContext)
    if method == POST || method == PUT || method == DELETE {
        return !frc.Exists() || (frc.IsDir() || (frc.IsFile() && frc.CanWrite(true))), "", req, cxt, 0, nil
    }
    return true, "", req, cxt, 0, nil
}

func (p *FileResource) Forbidden(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return false, req, cxt, 0, nil
}

func (p *FileResource) AllowMissingPost(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return true, req, cxt, 0, nil
}

func (p *FileResource) MalformedRequest(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return false, req, cxt, 0, nil
}

func (p *FileResource) URITooLong(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return len(req.URL().Path) > 4096, req, cxt, 0, nil
}

func (p *FileResource) DeleteResource(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    if !frc.Exists() {
        return true, req, cxt, 0, nil
    }
    path := frc.FullPath()
    var err os.Error
    if frc.IsFile() {
        err = os.Remove(path)
    } else if frc.IsDir() {
        err = os.RemoveAll(path)
    }
    if err == nil {
        return true, req, cxt, 0, nil
    }
    return false, req, cxt, http.StatusInternalServerError, err
}

/*
func (p *FileResource) DeleteCompleted(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
  return true, req, cxt, 0, nil
}
*/

func (p *FileResource) PostIsCreate(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return true, req, cxt, 0, nil
}

func (p *FileResource) CreatePath(req Request, cxt Context) (string, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    if frc.IsDir() {
        newPath := filepath.Join(frc.FullPath(), string(rand.Int63()))
        frc2 := NewFileResourceContextWithPath(newPath)
        for frc2.Exists() {
            newPath = filepath.Join(frc.FullPath(), string(rand.Int63()))
            frc2 = NewFileResourceContextWithPath(newPath)
        }
        frc = frc2
    } else if frc.Exists() {
        p := frc.FullPath()
        dir, tail := path.Split(p)
        ext := path.Ext(tail)
        basename := tail
        uniquify := time.UTC().Format(".20060102.150405")
        if len(ext) > 0 {
            basename = tail[:len(tail)-len(ext)] + uniquify
            frc.SetFullPath(path.Join(dir, basename+ext))
            for counter := 1; frc.Exists(); counter++ {
                frc.SetFullPath(path.Join(dir, basename+"."+strconv.Itoa(counter)+ext))
            }
        } else {
            basename = basename + uniquify
            frc.SetFullPath(path.Join(dir, basename))
            for counter := 1; frc.Exists(); counter++ {
                frc.SetFullPath(path.Join(dir, basename+"."+strconv.Itoa(counter)))
            }
        }
    }
    log.Print("[FileResource]: Will use path ", frc.FullPath())
    return frc.FullPath(), req, frc, 0, nil
}

func (p *FileResource) ProcessPost(req Request, cxt Context) (Request, Context, int, http.Header, io.WriterTo, os.Error) {
    // TODO handle POST
    mths, req, cxt, code, err := p.ContentTypesAccepted(req, cxt)
    if len(mths) > 0 {
        httpCode, httpHeaders, writerTo := mths[0].MediaTypeHandleInputFrom(req, cxt)
        return req, cxt, httpCode, httpHeaders, writerTo, nil
    }
    return req, cxt, code, nil, nil, err
}

func (p *FileResource) ContentTypesProvided(req Request, cxt Context) ([]MediaTypeHandler, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    var arr []MediaTypeHandler
    if frc.IsDir() {
        arr = []MediaTypeHandler{NewJsonDirectoryListing(frc.FullPath(), req.URL().Path), NewHtmlDirectoryListing(frc.FullPath(), req.URL().Path)}
    } else if frc.HasMultipleResources() {
        dir, _ := path.Split(frc.FullPath())
        filenames := frc.MultipleResourceNames()
        arr = make([]MediaTypeHandler, len(filenames))
        for i, filename := range filenames {
            extension := filepath.Ext(filename)
            mediaType := mime.TypeByExtension(extension)
            if len(mediaType) == 0 {
                // default to text/plain
                mediaType = MIME_TYPE_TEXT_PLAIN
            }
            fullFilename := path.Join(dir, filename)
            tempFrc := NewFileResourceContextWithPath(fullFilename)
            arr[i] = NewPassThroughMediaTypeHandler(mediaType, tempFrc, tempFrc.Len(), tempFrc.LastModified())
        }
    } else {
        extension := filepath.Ext(frc.FullPath())
        mediaType := mime.TypeByExtension(extension)
        if len(mediaType) == 0 {
            // default to text/plain
            mediaType = MIME_TYPE_TEXT_PLAIN
        }
        arr = []MediaTypeHandler{NewPassThroughMediaTypeHandler(mediaType, frc, frc.Len(), frc.LastModified())}
    }
    return arr, req, cxt, 0, nil
}

func (p *FileResource) ContentTypesAccepted(req Request, cxt Context) ([]MediaTypeInputHandler, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    mediaType := req.Header().Get("Content-Type")
    if len(mediaType) == 0 {
        extension := path.Ext(frc.FullPath())
        mediaType := mime.TypeByExtension(extension)
        if len(mediaType) == 0 {
            // default to text/plain
            mediaType = MIME_TYPE_TEXT_PLAIN
        }
    }
    knownContentLengthStr := req.Header().Get("Content-Length")
    knownContentLength := int64(-1)
    if len(knownContentLengthStr) > 0 {
        var err os.Error
        knownContentLength, err = strconv.Atoi64(knownContentLengthStr)
        if err != nil {
            knownContentLength = -1
        }
    }
    arr := []MediaTypeInputHandler{NewPassThroughMediaTypeInputHandler(mediaType, "", "", frc.FullPath(), path.Join(p.urlPathPrefix, path.Base(frc.FullPath())), false, knownContentLength, req.Body())}
    return arr, req, cxt, 0, nil
}

/*
func (p *FileResource) IsLanguageAvailable(languages []string, req Request, cxt Context) (bool, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) CharsetsProvided(charsets []string, req Request, cxt Context) ([]CharsetHandler, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) EncodingsProvided(encodings []string, req Request, cxt Context) ([]EncodingHandler, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) Variances(req Request, cxt Context) ([]string, Request, Context, int, os.Error) {

}
*/

func (p *FileResource) IsConflict(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    return frc.Exists() && !frc.IsFile(), req, cxt, 0, nil
}

func (p *FileResource) MultipleChoices(req Request, cxt Context) (bool, http.Header, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    if frc.HasMultipleResources() {
        headers := make(http.Header)
        headers.Set("Vary", "negotiate,accept")
        headers.Set("TCN", "choice")
        headers.Set("Accept-Ranges", "bytes")
        filenames := frc.MultipleResourceNames()
        contentTypeToFilename := make(map[string]string, len(filenames))
        contentTypes := make([]string, len(filenames))
        for i, filename := range filenames {
            extension := filepath.Ext(filename)
            mediaType := mime.TypeByExtension(extension)
            if len(mediaType) == 0 {
                // default to text/plain
                mediaType = MIME_TYPE_TEXT_PLAIN
            }
            contentTypeToFilename[mediaType] = filename
            contentTypes[i] = mediaType
        }
        finalContentType := chooseMediaTypeDefault(contentTypes, req.Header().Get("Accept"), contentTypes[0])
        finalFilename := contentTypeToFilename[finalContentType]
        headers.Set("Content-Type", finalContentType)
        headers.Set("Content-Location", finalFilename)
        dir, _ := path.Split(frc.FullPath())
        cxt = NewFileResourceContextWithPath(path.Join(dir, finalFilename))
        return true, headers, req, cxt, 0, nil
    }
    return false, nil, req, cxt, 0, nil
}

/*
func (p *FileResource) PreviouslyExisted(req Request, cxt Context) (bool, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) MovedPermanently(req Request, cxt Context) (string, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) MovedTemporarily(req Request, cxt Context) (string, Request, Context, int, os.Error) {

}
*/

func (p *FileResource) LastModified(req Request, cxt Context) (*time.Time, Request, Context, int, os.Error) {
    frc := cxt.(FileResourceContext)
    return frc.LastModified(), req, cxt, 0, nil
}
/*
func (p *FileResource) Expires(req Request, cxt Context) (*time.Time, Request, Context, int, os.Error) {

}
*/
/*
func (p *FileResource) GenerateETag(req Request, cxt Context) (string, Request, Context, int, os.Error) {

}
*/

func (p *FileResource) FinishRequest(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    if frc, ok := cxt.(FileResourceContext); ok {
        frc.Close()
    }
    return true, req, cxt, 0, nil
}

func (p *FileResource) ResponseIsRedirect(req Request, cxt Context) (bool, Request, Context, int, os.Error) {
    return false, req, cxt, 0, nil
}

func (p *FileResource) HasRespBody(req Request, cxt Context) bool {
    method := req.Method()
    if method == HEAD || method == PUT || method == DELETE {
        return false
    }
    return true
}
