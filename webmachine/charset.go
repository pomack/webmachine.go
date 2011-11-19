package webmachine

import (
    "io"
)

type StandardCharsetHandler struct {
    charset string
}

func NewStandardCharsetHandler(charset string) *StandardCharsetHandler {
    return &StandardCharsetHandler{charset: charset}
}

func (p *StandardCharsetHandler) Charset() string {
    return p.charset
}

func (p *StandardCharsetHandler) CharsetConverter(req Request, cxt Context, reader io.Reader) io.Reader {
    return reader
}

func (p *StandardCharsetHandler) String() string {
    return "NewStandardCharsetHandler(\"" + p.charset + "\")"
}
