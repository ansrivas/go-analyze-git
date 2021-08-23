// MIT License
//
// Copyright (c) 2021 Ankur Srivastava
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package fileops

import (
	"bufio"
	"fmt"
	"os"
)

type FileOps struct {
	// Using bufio.Scanner comes with a quirk, buffer size
	// is fixed, therefore this can be used to modify the bufsize
	BufSize int
}

func New() *FileOps {
	return &FileOps{
		BufSize: 64 * 1024,
	}
}

func NewWithBufSize(bufSize int) *FileOps {
	return &FileOps{
		BufSize: bufSize,
	}
}

func (f *FileOps) readFileStreaming(fname string, outputChan chan<- string, splitFunc bufio.SplitFunc) error {

	defer close(outputChan)

	fileHandle, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("failed to open file at path %s with error %v", fname, err.Error())
	}
	defer fileHandle.Close()

	fileReader := bufio.NewScanner(fileHandle)
	fileReader.Split(splitFunc)
	buf := make([]byte, f.BufSize)
	fileReader.Buffer(buf, f.BufSize)

	for fileReader.Scan() {
		outputChan <- fileReader.Text()
	}

	if err := fileReader.Err(); err != nil {
		return fmt.Errorf("shouldn't see an error scanning a string %s", err.Error())
	}
	return nil
}

// ReadFileStreaming reads a given file path
// and publishes it on a output channel.
//
// NOTE: Once the file is read, the channel should be closed
// from inside the function. DO NOT close it from outside.
//
// **This function is intended to be used as a Go-Routine**
func (f *FileOps) ReadFileStreaming(fname string, outputChan chan<- string, splitFunc bufio.SplitFunc) error {
	return f.readFileStreaming(fname, outputChan, splitFunc)
}
