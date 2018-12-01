package main

import (
	"bytes"
	"io"
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
