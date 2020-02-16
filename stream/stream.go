package stream

import (
	"io"
)

const bufferSize = 1024 * 4

type Streamer interface {
	Into(dst io.WriteCloser)
	Open() error
	Reset()
	Close() error
}

func NewReader(b []byte) *Reader {
	return &Reader{
		data: b,
	}
}

type Reader struct {
	data []byte
	cur  int
}

func (r *Reader) Read(b []byte) (n int, err error) {
	l := len(r.data)

	if r.cur >= l {
		return 0, io.EOF
	}

	to := r.cur + len(b)

	if to >= l {
		to = l
	}

	copy(b, r.data[r.cur:to])

	n = to - r.cur
	r.cur += n

	return
}

func (r *Reader) Close() error {
	*r = Reader{}

	return nil
}

func NewWriter(b []byte) *Writer {
	return &Writer{
		data: b,
	}
}

type Writer struct {
	data []byte
	cur  int
}

func (w *Writer) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, io.ErrShortWrite
	}

	w.data = append(w.data, b...)
	n = len(b)

	return
}

func (w *Writer) Close() error {
	*w = Writer{}

	return nil
}

func NewStream(src io.ReadCloser) *dsssStream {
	return &dsssStream{
		src: src,
	}
}

type dsssStream struct {
	src  io.ReadCloser
	dst  io.WriteCloser
	buff []byte
}

func (s *dsssStream) Into(dst io.WriteCloser) {
	s.dst = dst
}

func (s *dsssStream) Open() error {
	var isEOF bool
	s.buff = make([]byte, bufferSize)

	for !isEOF {
		n, err := s.src.Read(s.buff)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		_, err = s.dst.Write(s.buff[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *dsssStream) Reset() {
	*s = dsssStream{}
}

func (s *dsssStream) Close() error {
	s.buff = nil
	if err := s.src.Close(); err != nil {
		return err
	}

	return s.dst.Close()
}
