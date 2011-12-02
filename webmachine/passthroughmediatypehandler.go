package webmachine

import (
    "http"
    "io"
    //"log"
    "strconv"
    "time"
    "os"
)

type PassThroughMediaTypeHandler struct {
    mediaType           string
    reader              io.ReadCloser
    numberOfBytes       int64
    lastModified        *time.Time
    statusCode          int
    writtenStatusHeader bool
}

func NewPassThroughMediaTypeHandler(mediaType string, reader io.ReadCloser, numberOfBytes int64, lastModified *time.Time) *PassThroughMediaTypeHandler {
    return &PassThroughMediaTypeHandler{
        mediaType:     mediaType,
        reader:        reader,
        numberOfBytes: numberOfBytes,
        lastModified:  lastModified,
    }
}

func (p *PassThroughMediaTypeHandler) MediaTypeOutput() string {
    return p.mediaType
}

func (p *PassThroughMediaTypeHandler) SetStatusCode(statusCode int) {
    p.statusCode = statusCode
}

func (p *PassThroughMediaTypeHandler) MediaTypeHandleOutputTo(req Request, cxt Context, writer io.Writer, resp ResponseWriter) {
    if !p.writtenStatusHeader {
        if p.statusCode <= 0 {
            p.statusCode = http.StatusOK
        }
        resp.WriteHeader(p.statusCode)
        p.writtenStatusHeader = true
    }
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
            resp.Header().Set("Content-Range", "bytes="+outRangeString)
            currentOffset := int64(0)
            for _, arange := range ranges {
                start := arange[0]
                end := arange[1]
                var err os.Error
                if currentOffset < start {
                    if seeker, ok := p.reader.(io.Seeker); ok {
                        currentOffset, err = seeker.Seek(start-currentOffset, 1)
                        if err != nil {
                            return
                        }
                    } else {
                        if start-currentOffset >= 32768 {
                            buf := make([]byte, 32768)
                            for ; currentOffset+32768 < start; currentOffset += 32768 {
                                if _, err = io.ReadFull(p.reader, buf); err != nil {
                                    return
                                }
                            }
                        }
                        if currentOffset < start {
                            buf := make([]byte, start-currentOffset)
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
                    written, err := io.Copyn(writer, p.reader, end-currentOffset)
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
    //log.Print("[PTMTH]: Writer ", writer, "\n[PTMTH]: Reader ", p.reader, "\n[PTMTH]: numBytes ", p.numberOfBytes, "\n[PTMTH]: currentOffset ", currentOffset, "\n")
    for currentOffset < int64(p.numberOfBytes) {
        bytesToSend := p.numberOfBytes - currentOffset
        data := make([]byte, bytesToSend+10000)
        numBytesRead, err := p.reader.Read(data[0:bytesToSend])
        currentOffset += int64(numBytesRead)
        if err != nil {
            return
        }
        //log.Print("[PTMTH]: About to write ", len(data[0:bytesToSend]), " bytes to the writer\n")
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
