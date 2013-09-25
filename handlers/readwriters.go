package handlers

import (
	"io"
)

type multiWriteCloser struct {
	writers     []io.WriteCloser
	multiWriter io.Writer
}

type multiReadCloser struct {
	readers     []io.ReadCloser
	multiReader io.Reader
}

type NopWriter struct {
	w io.Writer
}

func newMultiWriterCloser(w ...io.WriteCloser) *multiWriteCloser {
	var mwc = new(multiWriteCloser)
	writers := make([]io.Writer, len(w))
	for index, wr := range w {
		writers[index] = io.Writer(wr)
	}
	mwc.multiWriter = io.MultiWriter(writers...)
	mwc.writers = w
	return mwc
}

func newMultiReaderCloser(r ...io.ReadCloser) *multiReadCloser {
	var mwr = new(multiReadCloser)
	readers := make([]io.Reader, len(r))
	for index, wr := range r {
		readers[index] = io.Reader(wr)
	}
	mwr.multiReader = io.MultiReader(readers...)
	mwr.readers = r
	return mwr
}

func newNopWriter(w io.Writer) *NopWriter {
	hw := &NopWriter{w: w}
	return hw
}

func (mwc *multiWriteCloser) Close() error {
	for _, w := range mwc.writers {
		w.Close()
	}
	return nil
}

func (mwr *multiReadCloser) Read(p []byte) (int, error) {
	return mwr.multiReader.Read(p)
}

func (mwr *multiReadCloser) Close() error {
	for _, r := range mwr.readers {
		r.Close()
	}
	return nil
}

func (mwc *multiWriteCloser) Write(p []byte) (int, error) {
	return mwc.multiWriter.Write(p)
}
func (np *NopWriter) Close() error {
	//NoOp
	return nil
}

func (hw *NopWriter) Write(p []byte) (int, error) {
	return hw.w.Write(p)
}
