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

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestReadFileStreaming(t *testing.T) {

	assert := assert.New(t)
	outputChan := make(chan string, 1)

	var g run.Group
	{
		g.Add(func() error {
			fileops := New()
			return fileops.ReadFileStreaming("testdata/01_data.csv", outputChan, bufio.ScanLines)
		}, func(err error) {
			if err != nil {
				log.Error().Msgf("The final goroutine actor was interrupted with: %v\n", err)
			}
		})

		g.Add(func() error {
			line1 := <-outputChan
			line2 := <-outputChan
			assert.Equal("header1,header2", line1)
			assert.Equal("test1,test2", line2)
			return nil
		}, func(err error) {
			// pass
		})
	}

	err := g.Run()
	assert.Nil(err, "No error should have been found")
}

func TestReadFileStreamingLimitedBuffer(t *testing.T) {
	assert := assert.New(t)

	outputChan := make(chan string, 1)
	var g run.Group
	{
		g.Add(func() error {
			fileops := New()
			return fileops.ReadFileStreaming("testdata/02_data_long.csv", outputChan, bufio.ScanLines)
		}, func(err error) {
			if err != nil {
				log.Error().Msgf("The final goroutine actor was interrupted with: %v\n", err)
			}
		})

		g.Add(func() error {
			for row := range outputChan {
				println(row)
			}
			return nil
		}, func(err error) {
			// pass
		})
	}
	err := g.Run()
	assert.NotNil(err)
}

func TestReadFileStreamingExtendedBuffer(t *testing.T) {
	assert := assert.New(t)
	outputChan := make(chan string, 1)

	var g run.Group
	{
		g.Add(func() error {
			fileops := NewWithBufSize(128 * 1024)
			return fileops.ReadFileStreaming("testdata/02_data_long.csv", outputChan, bufio.ScanLines)
		}, func(err error) {
			if err != nil {
				log.Error().Msgf("The final goroutine actor was interrupted with: %v\n", err)
			}
		})

		g.Add(func() error {
			for row := range outputChan {
				println(row)
			}
			return nil
		}, func(err error) {
			// pass
		})
	}
	err := g.Run()
	assert.Nil(err)
}
