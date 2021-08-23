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
)

// GenericDict represents a generic map represented as key, value
type GenericDict struct {
	Key   string `json:"Key"`
	Value int    `json:"Value"`
}

type GenericDictHeap []GenericDict

func (g GenericDictHeap) Len() int           { return len(g) }
func (g GenericDictHeap) Less(i, j int) bool { return g[i].Value < g[j].Value }
func (g GenericDictHeap) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

func (g GenericDictHeap) Take(count int) GenericDictHeap {
	if count < len(g) {
		return g[0:count]
	}
	return g
}

func (g *GenericDictHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// Taken from https://www.tugberkugurlu.com/archive/usage-of-the-heap-data-structure-in-go-golang-with-examples
	*g = append(*g, x.(GenericDict))
}

func (g *GenericDictHeap) Pop() interface{} {
	// Taken from https://www.tugberkugurlu.com/archive/usage-of-the-heap-data-structure-in-go-golang-with-examples
	old := *g
	n := len(old)
	x := old[n-1]
	*g = old[0 : n-1]
	return x
}

// Write the generated json to output
func (g GenericDictHeap) ToJson(out io.Writer) error {
	payload, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return err
	}
	_, err = out.Write(payload)
	return err
}
