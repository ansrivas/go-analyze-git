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

import "github.com/urfave/cli/v2"

// Repository struct is responsible for all the operations on
// repositories
type Repository struct{}

// New returns a new instance of repository
func New() *Repository {
	return &Repository{}
}

// topKReposByCommits returns Top K repositories by
// the amount of commits pushed
func (r *Repository) topKReposByCommits(num int) error {
	return nil
}

// topKReposByWatchEvents returns Top K repositories
// sorted by watch events
func (r *Repository) topKReposByWatchEvents(num int) error {
	return nil
}

func (r *Repository) CmdTopKReposByCommits() *cli.Command {
	return &cli.Command{
		Name:    "topk-by-commits",
		Aliases: []string{"tc"},
		Usage:   "Top K repositories by the amount of commits pushed",
		Action: func(*cli.Context) error {
			k := 10
			return r.topKReposByCommits(k)
		},
	}
}

func (r *Repository) CmdTopKReposByWatchEvents() *cli.Command {
	return &cli.Command{
		Name:    "topk-by-watchevents",
		Aliases: []string{"tw"},
		Usage:   "Top K repositories sorted by watch events",
		Action: func(*cli.Context) error {
			k := 10
			return r.topKReposByWatchEvents(k)
		},
	}
}
