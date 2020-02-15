package stream

import "io"

type testStruct struct {
	body []byte
	cur  int
}

func (f *testStruct) Read(p []byte) (n int, err error) {
	l := len(f.body)

	if f.cur >= l {
		return 0, io.EOF
	}

	to := f.cur + len(p)

	if to >= l {
		to = l
	}

	copy(p, f.body[f.cur:to])

	n = to - f.cur
	f.cur += n

	return
}

func (f *testStruct) Write(p []byte) (n int, err error) {
	f.body = append(f.body, p...)

	return len(p), nil
}

func (f *testStruct) Close() error {
	*f = testStruct{
		body: nil,
		cur:  0,
	}

	return nil
}
