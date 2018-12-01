package main

import (
	"bytes"
	"io"
	"os"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {

			break
		}

		if err != nil {
			return count, err
		}

	}

	return count, nil

}

func writeLine(file *os.File, host string) error {

	mutex.Lock()
	defer mutex.Unlock()

	_, err := file.WriteString(host + "\n")

	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) (file *os.File, exists bool, err error) {

	if _, err = os.Stat(output); os.IsNotExist(err) {
		return nil, false, nil
	} else {
		file, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		return file, true, err
	}
}

func createFile(filename string) (outputFile *os.File, err error) {

	file, err := os.Create(filename)

	if err != nil {
		return
	}

	defer file.Close()

	outputFile, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return
	}

	return
}

func initOutputFile(output string) (outputFile *os.File, err error) {

	outputFile, exists, err = fileExists(output)

	if err != nil {
		return
	}

	if !exists {
		outputFile, err = createFile(output)

		if err != nil {
			return
		}
	} else {
		outputFile, err = os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

		if err != nil {
			return
		}
	}

	return
}
