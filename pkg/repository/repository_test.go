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

package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
	"gitlab.com/ansrivas/go-analyze-git/pkg/events"
)

func TestTopKReposByEvents(t *testing.T) {
	assert := assert.New(t)

	repos := New()
	count := 2
	event := events.Watch
	eventsFile := "testdata/events.csv"
	reposFile := "testdata/repos.csv"
	cache, err := repos.topKReposByEvents(count, event, eventsFile, reposFile)

	expectedCache := make(map[string]int)
	expectedCache["testrepo2"] = 3
	expectedCache["testrepo1"] = 2
	expected := utils.GenericIntDictListFromMap(expectedCache)
	expected.SortByValue()

	assert.Equal(cache, expected)
	assert.Nil(err)
}
