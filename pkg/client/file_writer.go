package client

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type fileWriter struct {
	dataWriter *bufio.Writer
	file       *os.File
}

func newFileWriter(routine int) (io.WriteCloser, error) {
	filename := fmt.Sprintf("output/file%d.txt", routine)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %w", err)
	}

	return &fileWriter{file: file, dataWriter: bufio.NewWriter(file)}, nil
}

func (f *fileWriter) Write(m []byte) (int, error) {
	n, err := f.dataWriter.Write(m)
	if err != nil {
		return n, err
	}

	if err = f.dataWriter.WriteByte('\n'); err != nil {
		return n, err
	}
	return n + 1, nil
}

func (f *fileWriter) Close() error {
	f.dataWriter.Flush()
	return f.file.Close()
}
