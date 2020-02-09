package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readFile(filePath string) (filename string, content []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	if stat.IsDir() {
		err = errors.New("it is a dir")
		return
	}

	filename = filepath.Base(filePath)
	content = make([]byte, stat.Size())
	_, err = file.Read(content)

	return
}

func writeFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func writeToHistoryFile(filename, id string) error {
	t := time.Now()
	uploadDate := fmt.Sprintf("%d-%02d-%02d at %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	record := "uploaded: " + uploadDate + "\t" + id + "\t" + filename + "\n"

	return writeFile(historyFile, []byte(record))
}

func deleteFromHistoryFile(s string) error {
	_, content, err := readFile(historyFile)
	if err != nil {
		return err
	}

	t := time.Now()
	deletedDate := fmt.Sprintf("%d-%02d-%02d at %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	newVal := s + "\t deleted: " + deletedDate + "\n"

	newContent := strings.Replace(string(content), s+"\n", newVal, -1)

	return ioutil.WriteFile(historyFile, []byte(newContent), 0666)
}
