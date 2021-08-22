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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFileStreaming(t *testing.T) {
	assert := assert.New(t)

	fileops := New()

	outputChan := make(chan string, 1)
	errChan := make(chan error)
	go fileops.ReadFileStreaming("testdata/01_data.csv", outputChan, errChan, bufio.ScanLines)

	line1 := <-outputChan
	line2 := <-outputChan
	assert.Equal("header1,header2", line1)
	assert.Equal("test1,test2", line2)

	// This should return as the channel was closed
	<-outputChan
	err := <-errChan
	assert.Nil(err, "No error should have been found")
}

func TestReadFileStreamingLimitedBuffer(t *testing.T) {
	assert := assert.New(t)

	fileops := New()
	outputChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go fileops.ReadFileStreaming("testdata/02_data_long.csv", outputChan, errChan, bufio.ScanLines)

	// Try to read from outputChan
	for row := range outputChan {
		println(row)
	}
	err := <-errChan
	assert.NotNil(err)
}

func TestReadFileStreamingExtendedBuffer(t *testing.T) {
	assert := assert.New(t)

	fileops := NewWithBufSize(124 * 1024)
	outputChan := make(chan string, 1)
	errChan := make(chan error, 1)
	go fileops.ReadFileStreaming("testdata/02_data_long.csv", outputChan, errChan, bufio.ScanLines)

	// Try to read from outputChan
	for row := range outputChan {
		println(row)
	}
	err := <-errChan
	assert.Nil(err)
}
