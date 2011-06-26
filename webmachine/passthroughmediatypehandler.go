package webmachine

import (
  "container/list"
  "http"
  "io"
  "log"
  "strconv"
  "strings"
  "time"
  "os"
)


func NewPassThroughMediaTypeHandler(mediaType string, reader io.ReadCloser, numberOfBytes int64, lastModified *time.Time) *PassThroughMediaTypeHandler {
  log.Print("Creating Pass Through Media Type Handler: ", mediaType, ", ", reader, ", ", numberOfBytes, ", ", lastModified, "\n")
  return &PassThroughMediaTypeHandler{
    mediaType: mediaType,
    reader: reader,
    numberOfBytes: numberOfBytes,
    lastModified: lastModified,
  }
}

func (p *PassThroughMediaTypeHandler) MediaType() string {
  return p.mediaType
}

func (p *PassThroughMediaTypeHandler) splitRangeHeaderString(rangeHeader string) ([][2]int64) {
  if len(rangeHeader) > 6 && rangeHeader[0:6] == "bytes=" {
    rangeStrings := strings.Split(rangeHeader[6:], ",", -1)
    ranges := make([][2]int64, len(rangeStrings))
    for i, rangeString := range rangeStrings {
      trimmedRangeString := strings.TrimSpace(rangeString)
      dashIndex := strings.Index(rangeString, "-")
      switch {
      case dashIndex < 0:
        // single byte, e.g. 507
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][1] = ranges[i][0] + 1
      case dashIndex == 0:
        // start from end, e.g -51
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][0] += p.numberOfBytes
        ranges[i][1] = p.numberOfBytes
      case dashIndex == len(trimmedRangeString):
        // byte to end, e.g. 9500-
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][1] = p.numberOfBytes
      default:
        // range, e.g. 400-500
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString[0:dashIndex])
        ranges[i][1], _ = strconv.Atoi64(trimmedRangeString[dashIndex:])
        ranges[i][1] += 1
      }
      if ranges[i][0] >= p.numberOfBytes {
        continue
      }
      if ranges[i][1] > p.numberOfBytes {
        ranges[i][1] = p.numberOfBytes
      }
      // TODO sorting and compression of byte ranges
    }
    // sort ranges in ascending order
    for i, arange := range ranges {
      for ; i > 0 && ranges[i-1][0] > arange[0]; i-- {
        ranges[i-1][0], ranges[i-1][1], ranges[i][0], ranges[i][1] = ranges[i][0], ranges[i][1], ranges[i-1][0], ranges[i-1][1]
      }
    }
    // perform range compression for non-canonical ranges
    l := list.New()
    lastRange := ranges[0]
    for i, arange := range ranges {
      if i == 0 || lastRange[1] < arange[0] {
        l.PushBack(arange)
        lastRange = arange
      } else if lastRange[1] >= arange[0] {
        if lastRange[1] < arange[1] {
          lastRange[1] = arange[1]
        }
      } else {
        l.PushBack(arange)
        lastRange = arange
      }
    }
    theranges := make([][2]int64, l.Len())
    for i, elem :=0, l.Front(); elem != nil; elem, i = elem.Next(), i + 1 {
      theranges[i] = elem.Value.([2]int64)
    }
    return theranges
  }
  theranges := make([][2]int64, 1)
  theranges[0][0] = 0
  theranges[0][1] = p.numberOfBytes
  return theranges
}

func (p *PassThroughMediaTypeHandler) OutputTo(req Request, cxt Context, writer io.Writer, resp http.ResponseWriter) {
  if req.Header().Get("Accept-Ranges") == "bytes" {
    rangeHeader := req.Header().Get("Range")
    if len(rangeHeader) > 6 && rangeHeader[0:6] == "bytes=" {
      ranges := p.splitRangeHeaderString(rangeHeader)
      outRangeString := "bytes="
      for i, arange := range ranges {
        if i > 0 {
          outRangeString += ","
        }
        outRangeString += strconv.Itoa64(arange[0]) + "-" + strconv.Itoa64(arange[1]-1)
      }
      outRangeString += "/" + strconv.Itoa64(p.numberOfBytes)
      resp.Header().Set("Content-Range", "bytes=" + outRangeString)
      currentOffset := int64(0)
      for _, arange := range ranges {
        start := arange[0]
        end := arange[1]
        var err os.Error
        if currentOffset < start {
          if seeker, ok := p.reader.(io.Seeker); ok {
            currentOffset, err = seeker.Seek(start - currentOffset, 1)
            if err != nil {
              return
            }
          } else {
            if start - currentOffset >= 32768 {
              buf := make([]byte, 32768)
              for ; currentOffset + 32768 < start; currentOffset += 32768 {
                if _, err = io.ReadFull(p.reader, buf); err != nil {
                  return
                }
              }
            }
            if currentOffset < start {
              buf := make([]byte, start - currentOffset)
              if _, err = io.ReadFull(p.reader, buf); err != nil {
                return
              }
            }
          }
        }
        if req.Method() == HEAD {
          return
        }
        for currentOffset < end {
          written, err := io.Copyn(writer, p.reader, end - currentOffset)
          currentOffset += written
          if err != nil {
            return
          }
        }
      }
      return
    }
  }
  if req.Method() == HEAD {
    return
  }
  currentOffset := int64(0)
  log.Print("[PTMTH]: Writer ", writer, "\n[PTMTH]: Reader ", p.reader, "\n[PTMTH]: numBytes ", p.numberOfBytes, "\n[PTMTH]: currentOffset ", currentOffset, "\n")
  for currentOffset < int64(p.numberOfBytes) {
    bytesToSend := p.numberOfBytes - currentOffset
    data := make([]byte, bytesToSend  + 10000)
    numBytesRead, err := p.reader.Read(data[0:bytesToSend])
    currentOffset += int64(numBytesRead)
    if err != nil {
      return
    }
    log.Print("[PTMTH]: About to write ", len(data[0:bytesToSend]), " bytes to the writer\n")
    _, err = writer.Write(data[0:bytesToSend])
    if err != nil {
      return
    }
    //written, err := io.Copyn(writer, p.reader, p.numberOfBytes - currentOffset)
    //if err != nil {
    //  return
    //}
    //currentOffset += int64(written)
  }
}

