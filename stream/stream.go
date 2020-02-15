package stream

import (
	"io"
)

const (
	messageSize = 1024 * 4
	contentType = "multipart/form-data"
)

type Streamer interface {
	Into(dst io.WriteCloser)
	Open() error
	Reset()
	Close() error
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
	s.buff = make([]byte, messageSize)

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
