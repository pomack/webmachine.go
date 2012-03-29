package webmachine

import (
    "encoding/json"
    "io"
    "log"
    "os"
    "path"
    "time"
)

type jsonDirectoryEntry struct {
    Filename     string `json:"filename"`
    Path         string `json:"path"`
    Size         int64  `json:"size"`
    IsDirectory  bool   `json:"is_directory"`
    LastModified string `json:"last_modified"`
}

type jsonDirectoryEntryResult struct {
    Status  string               `json:"status"`
    Message string               `json:"message"`
    Path    string               `json:"path"`
    Result  []jsonDirectoryEntry `json:"result"`
}

type JsonDirectoryListing struct {
    fullPath string
    urlPath  string
    file     *os.File
}

type htmlDirectoryEntry struct {
    Filename     string `json:"filename"`
    Path         string `json:"path"`
    Size         int64  `json:"size"`
    IsDirectory  bool   `json:"is_directory"`
    LastModified string `json:"last_modified"`
}

type htmlDirectoryEntryResult struct {
    Status       string               `json:"status"`
    Tail         string               `json:"tail"`
    Path         string               `json:"path"`
    Message      string               `json:"message"`
    LastModified string               `json:"last_modified"`
    Result       []htmlDirectoryEntry `json:"result"`
}

type HtmlDirectoryListing struct {
    fullPath string
    urlPath  string
    file     *os.File
}

func NewJsonDirectoryListing(fullPath string, urlPath string) *JsonDirectoryListing {
    return &JsonDirectoryListing{fullPath: fullPath, urlPath: urlPath}
}

func (p *JsonDirectoryListing) MediaTypeOutput() string {
    return MIME_TYPE_JSON
}

func (p *JsonDirectoryListing) MediaTypeHandleOutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
    result := new(jsonDirectoryEntryResult)
    result.Path = p.urlPath
    var err error
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
            result.Message = err.Error()
            result.Result = make([]jsonDirectoryEntry, 0)
            encoder := json.NewEncoder(writer)
            encoder.Encode(result)
            return
        }
    }
    fileInfos, err := p.file.Readdir(-1)
    if err != nil {
        result.Status = "error"
        result.Message = err.Error()
        result.Result = make([]jsonDirectoryEntry, 0)
        encoder := json.NewEncoder(writer)
        encoder.Encode(result)
        return
    }
    entries := make([]jsonDirectoryEntry, len(fileInfos))
    for i, fileInfo := range fileInfos {
        entries[i].Filename = fileInfo.Name()
        entries[i].Path = path.Join(p.urlPath, fileInfo.Name())
        entries[i].Size = fileInfo.Size()
        entries[i].IsDirectory = fileInfo.IsDir()
        if fileInfo.IsDir() {
            entries[i].IsDirectory = true
            entries[i].Size = 0
        } else {
            entries[i].IsDirectory = false
            entries[i].Size = fileInfo.Size()
        }
        entries[i].LastModified = fileInfo.ModTime().UTC().Format(time.RFC3339)
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

func (p *HtmlDirectoryListing) MediaTypeOutput() string {
    return MIME_TYPE_HTML
}

func (p *HtmlDirectoryListing) MediaTypeHandleOutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
    result := new(htmlDirectoryEntryResult)
    result.Path = p.urlPath
    result.Tail = path.Base(p.urlPath)
    var err error
    defer func() {
        if p.file != nil {
            p.file.Close()
            p.file = nil
        }
    }()
    if p.file == nil {
        p.file, err = os.Open(p.fullPath)
        if err != nil {
            result.Message = err.Error()
            result.Result = make([]htmlDirectoryEntry, 0)
            HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
            return
        }
    }
    fileInfos, err := p.file.Readdir(-1)
    if err != nil {
        result.Message = err.Error()
        result.Result = make([]htmlDirectoryEntry, 0)
        HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
        return
    }
    entries := make([]htmlDirectoryEntry, len(fileInfos))
    for i, fileInfo := range fileInfos {
        entries[i].Filename = fileInfo.Name()
        entries[i].Path = path.Join(p.urlPath, fileInfo.Name())
        entries[i].Size = fileInfo.Size()
        entries[i].IsDirectory = fileInfo.IsDir()
        if fileInfo.IsDir() {
            entries[i].IsDirectory = true
            entries[i].Size = 0
        } else {
            entries[i].IsDirectory = false
            entries[i].Size = fileInfo.Size()
        }
        entries[i].LastModified = fileInfo.ModTime().UTC().Format(time.ANSIC)
    }
    dirInfo, _ := p.file.Stat()
    if dirInfo != nil {
        result.LastModified = dirInfo.ModTime().UTC().Format(time.ANSIC)
    }
    result.Status = "success"
    result.Message = ""
    result.Result = entries
    log.Printf("Executing Success with result\n  %#v", result)
    HTML_DIRECTORY_LISTING_SUCCESS_TEMPLATE.ExecuteTemplate(writer, "directory_listing_success", result)
    return
}
