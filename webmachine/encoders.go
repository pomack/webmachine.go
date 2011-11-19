package webmachine

import (
  "compress/gzip"
  "compress/flate"
  "compress/lzw"
  "http"
  "io"
)

type identityEncoding struct {}

type compressEncoding struct {}

type gzipEncoding struct {}

type deflateEncoding struct {}

type chunkedEncoding struct {}

func NewIdentityEncoder() EncodingHandler {
  return new(identityEncoding)
}

func (p *identityEncoding) Encoding() string {
  return ENCODING_IDENTITY
}

func (p *identityEncoding) Encoder(req Request, cxt Context, writer io.Writer) (io.Writer) {
  return writer
}

func (p *identityEncoding) Decoder(req Request, cxt Context, reader io.Reader) (io.Reader) {
  return reader
}

func (p *identityEncoding) String() string {
  return ENCODING_IDENTITY
}


func NewCompressEncoder() EncodingHandler {
  return new(compressEncoding)
}

func (p *compressEncoding) Encoding() string {
  return ENCODING_COMPRESS
}

func (p *compressEncoding) Encoder(req Request, cxt Context, writer io.Writer) (io.Writer) {
  return lzw.NewWriter(writer, lzw.LSB, 8)
}

func (p *compressEncoding) Decoder(req Request, cxt Context, reader io.Reader) (io.Reader) {
  return lzw.NewReader(reader, lzw.LSB, 8)
}

func (p *compressEncoding) String() string {
  return ENCODING_COMPRESS
}

func NewGZipEncoder() EncodingHandler {
  return new(gzipEncoding)
}

func (p *gzipEncoding) Encoding() string {
  return ENCODING_GZIP
}

func (p *gzipEncoding) Encoder(req Request, cxt Context, writer io.Writer) (io.Writer) {
  w, _ := gzip.NewWriter(writer)
  return w
}

func (p *gzipEncoding) Decoder(req Request, cxt Context, reader io.Reader) (io.Reader) {
  r, _ := gzip.NewReader(reader)
  return r
}

func (p *gzipEncoding) String() string {
  return ENCODING_GZIP
}

func NewDeflateEncoder() EncodingHandler {
  return new(deflateEncoding)
}

func (p *deflateEncoding) Encoding() string {
  return ENCODING_DEFLATE
}

func (p *deflateEncoding) Encoder(req Request, cxt Context, writer io.Writer) (io.Writer) {
  w := flate.NewWriter(writer, flate.DefaultCompression)
  return w
}

func (p *deflateEncoding) Decoder(req Request, cxt Context, reader io.Reader) (io.Reader) {
  return flate.NewReader(reader)
}

func (p *deflateEncoding) String() string {
  return ENCODING_DEFLATE
}

func NewChunkedEncoder() EncodingHandler {
  return new(chunkedEncoding)
}

func (p *chunkedEncoding) Encoding() string {
  return ENCODING_CHUNKED
}

func (p *chunkedEncoding) Encoder(req Request, cxt Context, writer io.Writer) (io.Writer) {
  return http.NewChunkedWriter(writer)
}

func (p *chunkedEncoding) Decoder(req Request, cxt Context, reader io.Reader) (io.Reader) {
  return flate.NewReader(reader)
}

func (p *chunkedEncoding) String() string {
  return ENCODING_CHUNKED
}



