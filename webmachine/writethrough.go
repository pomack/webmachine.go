package webmachine

import (
    "io"
    "os"
)

func NewWriteThrough(from, to io.Writer) *WriteThrough {
    return &WriteThrough{from: from, to: to}
}

func (p *WriteThrough) Write(data []byte) (int, os.Error) {
    if len(data) == 0 {
        return 0, nil
    }
    return p.to.Write(data)
}

func (p *WriteThrough) Flush() os.Error {
    if flusher, ok := p.to.(Flusher); ok && flusher != nil {
        return flusher.Flush()
    }
    return nil
}

func (p *WriteThrough) Close() os.Error {
    if closer, ok := p.to.(io.Closer); ok && closer != nil {
        return closer.Close()
    }
    return nil
}
