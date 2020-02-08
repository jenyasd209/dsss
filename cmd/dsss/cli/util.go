package cli

import (
	"errors"
	"os"
	"strings"
)

func readFile(filePath string) (filename string, content []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}

	stat, err := file.Stat()
	if err != nil {
		return
	}

	if stat.IsDir() {
		err = errors.New("it is a dir")
		return
	}

	filename = file.Name()
	content = make([]byte, stat.Size())
	_, err = file.Read(content)

	return
}

func writeFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func newFileId(id string) error {
	return writeFile(idFile, []byte(id+"\n"))
}

func deleteFromFile(s string) error {
	_, content, err := readFile(idFile)
	if err != nil {
		return err
	}

	s = strings.Trim(s, string(content))

	return writeFile(idFile, []byte(s))
}
