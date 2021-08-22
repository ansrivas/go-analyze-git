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
package utils

import (
	"encoding/json"
	"io"
	"sort"
)

// GenericIntDict represents a generic map represented as key, value
type GenericIntDict struct {
	// UserID represents the readable user-id
	Key string `json:"Key"`
	// Count represents the count of Commits and PRs
	Value int `json:"Value"`
}

type GenericIntDictList []GenericIntDict

func (g GenericIntDictList) Len() int                          { return len(g) }
func (g GenericIntDictList) Less(i, j int) bool                { return g[i].Value < g[j].Value }
func (g GenericIntDictList) Swap(i, j int)                     { g[i], g[j] = g[j], g[i] }
func (g GenericIntDictList) Take(count int) GenericIntDictList { return g[0:count] }

// Write the generated json to output
func (g GenericIntDictList) ToJson(out io.Writer) error {
	payload, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return err
	}
	_, err = out.Write(payload)
	return err
}

// Create a GenericIntDictList from a map
func GenericIntDictListFromMap(input map[string]int) GenericIntDictList {
	inner := []GenericIntDict{}
	for k, v := range input {
		inner = append(inner, GenericIntDict{
			Key:   k,
			Value: v,
		})
	}
	return GenericIntDictList(inner)
}

// Sort the GenericIntDictList by its value
func (g GenericIntDictList) SortByValue() {
	sort.Sort(sort.Reverse(g))
}
