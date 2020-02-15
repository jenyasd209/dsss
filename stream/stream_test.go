package stream

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDsssStream_Fake(t *testing.T) {
	f1 := &testStruct{body: []byte("fake byte")}
	f2 := &testStruct{}

	stream := NewStream(f1)
	defer stream.Close()

	stream.Into(f2)

	err := stream.Open()
	assert.Nil(t, err, err)
	assert.Equal(t, f1.body, f2.body)
}

func TestDsssStream_File(t *testing.T) {
	firstPath := "file1.txt"
	secondPath := "file2.txt"

	f1, err := os.OpenFile(firstPath, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	require.Nil(t, err, err)
	defer f1.Close()

	f2, err := os.OpenFile(secondPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	require.Nil(t, err, err)
	defer f2.Close()

	stream := NewStream(f1)
	stream.Into(f2)

	err = stream.Open()
	assert.Nil(t, err, err)

	err = stream.Close()
	assert.Nil(t, err, err)

	f1, err = os.OpenFile(firstPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	require.Nil(t, err, err)
	defer f1.Close()

	c1, err := read(f1)
	assert.Equal(t, io.EOF, err)

	f2, err = os.OpenFile(secondPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	require.Nil(t, err, err)
	defer f2.Close()

	c2, err := read(f2)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, c1, c2)

	err = os.Remove(f1.Name())
	assert.Nil(t, err, err)

	err = os.Remove(f2.Name())
	assert.Nil(t, err, err)
}

func read(f io.Reader) (content []byte, err error) {
	buf := make([]byte, 1)
	for {
		_, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		content = append(content, buf...)
	}

	return content, io.EOF
}
