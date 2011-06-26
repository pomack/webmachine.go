package webmachine

import (
  "io"
  "json"
  "os"
  "path"
  "time"
)

type jsonDirectoryEntry struct {
  Filename string "filename"
  Path string "path"
  Size int64 "size"
  IsDirectory bool "is_directory"
  LastModified string "last_modified"
}

type jsonDirectoryEntryResult struct {
  Result string "result"
  Message string "message"
  Path string "path"
  Entries []jsonDirectoryEntry "entries"
}

type JsonDirectoryListing struct {
  fullPath string
  urlPath string
  file *os.File
}

type htmlDirectoryEntry struct {
  Filename string "filename"
  Path string "path"
  Size int64 "size"
  IsDirectory bool "is_directory"
  LastModified string "last_modified"
}

type htmlDirectoryEntryResult struct {
  Tail string "tail"
  Path string "path"
  Message string "message"
  LastModified string "last_modified"
  Entries []htmlDirectoryEntry "entries"
}

type HtmlDirectoryListing struct {
  fullPath string
  urlPath string
  file *os.File
}

func NewJsonDirectoryListing(fullPath string, urlPath string) *JsonDirectoryListing {
  return &JsonDirectoryListing{fullPath: fullPath, urlPath: urlPath}
}

func (p *JsonDirectoryListing) MediaType() string {
  return "application/json"
}

func (p *JsonDirectoryListing) OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
  result := new(jsonDirectoryEntryResult)
  result.Path = p.urlPath
  var err os.Error
  if p.file == nil {
    p.file, err = os.Open(p.fullPath)
    if err != nil {
      if p.file != nil {
        p.file.Close()
        p.file = nil
      }
      result.Result = "error"
      result.Message = err.String()
      result.Entries = make([]jsonDirectoryEntry, 0)
      encoder := json.NewEncoder(writer)
      encoder.Encode(result)
      return
    }
  }
  fileInfos, err := p.file.Readdir(-1)
  if err != nil {
    result.Result = "error"
    result.Message = err.String()
    result.Entries = make([]jsonDirectoryEntry, 0)
    encoder := json.NewEncoder(writer)
    encoder.Encode(result)
    p.file.Close()
    p.file = nil
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
  p.file.Close()
  p.file = nil
  result.Result = "success"
  result.Message = ""
  result.Entries = entries
  encoder := json.NewEncoder(writer)
  encoder.Encode(result)
}


func NewHtmlDirectoryListing(fullPath string, urlPath string) *HtmlDirectoryListing {
  return &HtmlDirectoryListing{fullPath: fullPath, urlPath: urlPath}
}

func (p *HtmlDirectoryListing) MediaType() string {
  return "text/html"
}

func (p *HtmlDirectoryListing) OutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
  result := new(htmlDirectoryEntryResult)
  result.Path = p.urlPath
  result.Tail = path.Base(p.urlPath)
  var err os.Error
  if p.file == nil {
    p.file, err = os.Open(p.fullPath)
    if err != nil {
      if p.file != nil {
        p.file.Close()
        p.file = nil
      }
      result.Message = err.String()
      result.Entries = make([]htmlDirectoryEntry, 0)
      HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
      return
    }
  }
  fileInfos, err := p.file.Readdir(-1)
  if err != nil {
    result.Message = err.String()
    result.Entries = make([]htmlDirectoryEntry, 0)
    HTML_DIRECTORY_LISTING_ERROR_TEMPLATE.Execute(writer, result)
    p.file.Close()
    p.file = nil
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
  p.file.Close()
  p.file = nil
  result.Message = ""
  result.Entries = entries
  HTML_DIRECTORY_LISTING_SUCCESS_TEMPLATE.Execute(writer, result)
}
