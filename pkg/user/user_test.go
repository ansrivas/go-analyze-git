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

package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
)

func TestTopKUsersByPRsAndCommits(t *testing.T) {
	assert := assert.New(t)

	user := New()
	count := 3
	eventsFile := "testdata/events.csv"
	commitsFile := "testdata/commits.csv"
	actorsFile := "testdata/actors.csv"
	cache, err := user.topKUsersByPRsAndCommits(count, actorsFile, eventsFile, commitsFile)

	expected := utils.GenericDictHeap{
		utils.GenericDict{Key: "Apexal", Value: 5},
		utils.GenericDict{Key: "anggi1234", Value: 4},
		utils.GenericDict{Key: "onosendi", Value: 3},
	}

	assert.Equal(cache, expected)
	assert.Nil(err)
}

func BenchmarkTopKUsersByPRsAndCommits(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	user := New()
	count := 10
	eventsFile := "../../data/events.csv"
	commitsFile := "../../data/commits.csv"
	actorsFile := "../../data/actors.csv"
	for i := 0; i < b.N; i++ {
		user.topKUsersByPRsAndCommits(count, actorsFile, eventsFile, commitsFile) //nolint
	}
}
