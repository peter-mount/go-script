package io

import "io"

type IO struct {
	Discard io.Writer
}

func newIO() *IO {
	return &IO{
		Discard: io.Discard,
	}
}

func (_ IO) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return io.Copy(dst, src)
}

func (_ IO) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	return io.CopyBuffer(dst, src, buf)
}

func (_ IO) CopyN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	return io.CopyN(dst, src, n)
}

func (_ IO) ReadAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func (_ IO) ReadAtLeast(r io.Reader, buf []byte, min int) (n int, err error) {
	return io.ReadAtLeast(r, buf, min)
}

func (_ IO) ReadFull(r io.Reader, buf []byte) (n int, err error) {
	return io.ReadFull(r, buf)
}

func (_ IO) WriteString(w io.Writer, s string) (n int, err error) {
	return io.WriteString(w, s)
}

func (_ IO) LimitReader(r io.Reader, n int64) io.Reader {
	return io.LimitReader(r, n)
}

func (_ IO) MultiReader(readers ...io.Reader) io.Reader {
	return io.MultiReader(readers...)
}

func (_ IO) MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

func (_ IO) Pipe() (*io.PipeReader, *io.PipeWriter) {
	return io.Pipe()
}

func (_ IO) TeeReader(r io.Reader, w io.Writer) io.Reader {
	return io.TeeReader(r, w)
}
