// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webmachine

import (
    "io"
    "os"
    "strconv"
)

// NewChunkedWriter returns a new writer that translates writes into HTTP
// "chunked" format before writing them to w.  Closing the returned writer
// sends the final 0-length chunk that marks the end of the stream.
func NewChunkedWriter(w io.Writer) io.WriteCloser {
    return &chunkedWriter{writer: w}
}

// Writing to ChunkedWriter translates to writing in HTTP chunked Transfer
// Encoding wire format to the underlying writer.
type chunkedWriter struct {
    writer io.Writer
}

// Write the contents of data as one chunk to writer.
// NOTE: Note that the corresponding chunk-writing procedure in Conn.Write has
// a bug since it does not check for success of io.WriteString
func (p *chunkedWriter) Write(data []byte) (n int, err os.Error) {

    // Don't send 0-length data. It looks like EOF for chunked encoding.
    if len(data) == 0 {
        return 0, nil
    }

    head := strconv.Itob(len(data), 16) + "\r\n"

    if _, err = io.WriteString(p.writer, head); err != nil {
        return 0, err
    }
    if n, err = p.writer.Write(data); err != nil {
        return
    }
    if n != len(data) {
        err = io.ErrShortWrite
        return
    }
    _, err = io.WriteString(p.writer, "\r\n")

    return
}

func (p *chunkedWriter) Close() os.Error {
    if p.writer != nil {
        var err2 os.Error
        _, err := io.WriteString(p.writer, "0\r\n")
        if closer, ok := p.writer.(io.Closer); ok {
            err2 = closer.Close()
        }
        p.writer = nil
        if err == nil {
            return err2
        }
        return err
    }
    return nil
}

// NewChunkedReader returns a new reader that translates reads from HTTP
// "chunked" format before writing them to w.  Closing the returned writer
// sends the final 0-length chunk that marks the end of the stream.
func NewChunkedReader(r io.Reader) io.ReadCloser {
    return &chunkedReader{reader: r}
}

// Reading from ChunkedReader translates to writing in HTTP chunked Transfer
// Encoding wire format from the underlying reader.
type chunkedReader struct {
    reader           io.Reader
    bytesLeftInChunk int64
    endOfSection     bool
}

// Read the contents of data as one chunk from reader.
// NOTE: Note that the corresponding chunk-writing procedure in Conn.Write has
// a bug since it does not check for success of io.WriteString
func (p *chunkedReader) Read(data []byte) (n int, err os.Error) {

    // Don't send 0-length data. It looks like EOF for chunked encoding.
    if len(data) == 0 {
        return 0, nil
    }

    if p.bytesLeftInChunk != 0 {
        if p.bytesLeftInChunk > int64(len(data)) {
            n, err := p.reader.Read(data)
            p.bytesLeftInChunk -= int64(n)
            return n, err
        }
        n, err := p.reader.Read(data[0 : len(data)-int(p.bytesLeftInChunk)])
        p.bytesLeftInChunk -= int64(n)
        if err != nil {
            return n, err
        }
        n2, err := p.Read(data[n:])
        return n + n2, err
    }

    oneByte := make([]byte, 1)
    nl := make([]byte, 2)
    line := make([]byte, 16)
    offset := 0
    if p.endOfSection {
        bytesRead, err := p.reader.Read(nl)
        if err != nil || bytesRead != len(nl) || string(bytesRead) != "\r\n" {
            return 0, err
        }
    } else {
        p.endOfSection = true
    }
    for {
        bytesRead, err := p.reader.Read(oneByte)
        if err != nil || bytesRead != 1 {
            return 0, err
        }
        if oneByte[0] == '\r' {
            // discard \n
            p.reader.Read(oneByte)
            break
        }
        line[offset] = oneByte[0]
        offset++
    }
    p.bytesLeftInChunk, err = strconv.Atoi64(string(line[0 : offset+1]))
    if err != nil {
        return 0, err
    }
    return p.Read(data)
}

func (p *chunkedReader) Close() os.Error {
    if p.reader != nil {
        if closer, ok := p.reader.(io.Closer); ok {
            err := closer.Close()
            p.reader = nil
            return err
        }
        p.reader = nil
    }
    return nil
}
