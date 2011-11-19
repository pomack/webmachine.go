package webmachine

import (
    "io"
    "json"
    "os"
    "path"
    "time"
)

type jsonDirectoryEntry struct {
    Filename     string "filename"
    Path         string "path"
    Size         int64  "size"
    IsDirectory  bool   "is_directory"
    LastModified string "last_modified"
}

type jsonDirectoryEntryResult struct {
    Status  string               "status"
    Message string               "message"
    Path    string               "path"
    Result  []jsonDirectoryEntry "result"
}

type JsonDirectoryListing struct {
    fullPath string
    urlPath  string
    file     *os.File
}

type htmlDirectoryEntry struct {
    Filename     string "filename"
    Path         string "path"
    Size         int64  "size"
    IsDirectory  bool   "is_directory"
    LastModified string "last_modified"
}

type htmlDirectoryEntryResult struct {
    Status       string               "status"
    Tail         string               "tail"
    Path         string               "path"
    Message      string               "message"
    LastModified string               "last_modified"
    Result       []htmlDirectoryEntry "result"
}

type HtmlDirectoryListing struct {
    fullPath string
    urlPath  string
    file     *os.File
}

func NewJsonDirectoryListing(fullPath string, urlPath string) *JsonDirectoryListing {
    return &JsonDirectoryListing{fullPath: fullPath, urlPath: urlPath}
}

func (p *JsonDirectoryListing) MediaType() string {
    return MIME_TYPE_JSON
}

func (p *JsonDirectoryListing) OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
    result := new(jsonDirectoryEntryResult)
    result.Path = p.urlPath
    var err os.Error
    defer func() {
        if p.file != nil {
            p.file.Close()
            p.file = nil
        }
    }()
    if p.file == nil {
        p.file, err = os.Open(p.fullPath)
        if err != nil {
            result.Status = "error"
            result.Message = err.String()
            result.Result = make([]jsonDirectoryEntry, 0)
            encoder := json.NewEncoder(writer)
            encoder.Encode(result)
            return
        }
    }
    fileInfos, err := p.file.Readdir(-1)
    if err != nil {
        result.Status = "error"
        result.Message = err.String()
        result.Result = make([]jsonDirectoryEntry, 0)
        encoder := json.NewEncoder(writer)
        encoder.Encode(result)
        return
    }
    entries := make([]jsonDirectoryEntry, len(fileInfos))
    for i, fileInfo := range fileInfos {
        entries[i].Filename = fileInfo.Name
        entries[i].Path = path.Join(p.urlPath, fileInfo.Name)
        entries[i].Size = fileInfo.Size
        entries[i].IsDirectory = fileInfo.IsDirectory()
        if fileInfo.IsDirectory() {
            entries[i].IsDirectory = true
            entries[i].Size = 0
        } else {
            entries[i].IsDirectory = false
            entries[i].Size = fileInfo.Size
        }
        entries[i].LastModified = time.SecondsToUTC(int64(fileInfo.Mtime_ns / 1e9)).Format(time.RFC3339)
    }
    result.Status = "success"
    result.Message = ""
    result.Result = entries
    encoder := json.NewEncoder(writer)
    encoder.Encode(result)
}

func NewHtmlDirectoryListing(fullPath string, urlPath string) *HtmlDirectoryListing {
    return &HtmlDirectoryListing{fullPath: fullPath, urlPath: urlPath}
}

func (p *HtmlDirectoryListing) MediaType() string {
    return MIME_TYPE_HTML
}

func (p *HtmlDirectoryListing) OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
    result := new(htmlDirectoryEntryResult)
    result.Path = p.urlPath
    result.Tail = path.Base(p.urlPath)
    var err os.Error
    defer func() {
        if p.file != nil {
            p.file.Close()
            p.file = nil
        }
    }()
    if p.file == nil {
        p.file, err = os.Open(p.fullPath)
        if err != nil {
            result.Message = err.String()
            result.Result = make([]htmlDirectoryEntry, 0)
            HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
            return
        }
    }
    fileInfos, err := p.file.Readdir(-1)
    if err != nil {
        result.Message = err.String()
        result.Result = make([]htmlDirectoryEntry, 0)
        HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
        return
    }
    entries := make([]htmlDirectoryEntry, len(fileInfos))
    for i, fileInfo := range fileInfos {
        entries[i].Filename = fileInfo.Name
        entries[i].Path = path.Join(p.urlPath, fileInfo.Name)
        entries[i].Size = fileInfo.Size
        entries[i].IsDirectory = fileInfo.IsDirectory()
        if fileInfo.IsDirectory() {
            entries[i].IsDirectory = true
            entries[i].Size = 0
        } else {
            entries[i].IsDirectory = false
            entries[i].Size = fileInfo.Size
        }
        entries[i].LastModified = time.SecondsToUTC(int64(fileInfo.Mtime_ns / 1e9)).Format(time.ANSIC)
    }
    dirInfo, _ := p.file.Stat()
    if dirInfo != nil {
        result.LastModified = time.SecondsToUTC(int64(dirInfo.Mtime_ns / 1e9)).Format(time.ANSIC)
    }
    result.Status = "success"
    result.Message = ""
    result.Result = entries
    HTML_DIRECTORY_LISTING_SUCCESS_TEMPLATE.Execute(writer, result)
}
