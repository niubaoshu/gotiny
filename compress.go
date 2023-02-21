package gotiny

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
)

var (
	Gzip = GzipPool{
		readers: sync.Pool{},
		writers: sync.Pool{},
		bufferPool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
)

// GzipPool manages a pool of gzip.Writer.
// The pool uses sync.Pool internally.
type GzipPool struct {
	readers    sync.Pool
	writers    sync.Pool
	bufferPool sync.Pool
}

// get a buffer from pool
func (pool *GzipPool) Getbuffer() *bytes.Buffer {
	b := pool.bufferPool.Get().(*bytes.Buffer)
	b.Reset()
	return b
}

// put back a buffer to the pool
func (pool *GzipPool) Putbuffer(b *bytes.Buffer) {
	pool.bufferPool.Put(b)
}

// GetReader returns gzip.Reader from the pool, or creates a new one
// if the pool is empty.
func (pool *GzipPool) GetReader(src io.Reader) (reader *gzip.Reader) {
	if r := pool.readers.Get(); r != nil {
		reader = r.(*gzip.Reader)
		reader.Reset(src)
	} else {
		reader, _ = gzip.NewReader(src)
	}
	return reader
}

// PutReader closes and returns a gzip.Reader to the pool
// so that it can be reused via GetReader.
func (pool *GzipPool) PutReader(reader *gzip.Reader) {
	reader.Close()
	pool.readers.Put(reader)
}

// GetWriter returns gzip.Writer from the pool, or creates a new one
// with gzip.BestCompression if the pool is empty.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.writers.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		return gzip.NewWriter(dst)
	}
	return writer
}

// PutWriter closes and returns a gzip.Writer to the pool
// so that it can be reused via GetWriter.
func (pool *GzipPool) PutWriter(writer *gzip.Writer) {
	writer.Close()
	pool.writers.Put(writer)
}

func Gziper(outPut *bytes.Buffer, inPut []byte) error {
	gz := Gzip.GetWriter(outPut)
	defer Gzip.PutWriter(gz)
	defer gz.Flush()
	gz.Write(inPut)
	return nil
}

func Gunziper(outPut *bytes.Buffer, input *bytes.Buffer) error {
	gz := Gzip.GetReader(input)
	defer Gzip.PutReader(gz)
	io.Copy(outPut, gz)
	return nil
}
