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
	cli "github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
)

// UsersByPRsAndCommits represents a list of GenericIntDict
type UsersByPRsAndCommits []utils.GenericDictHeap

// User struct defines all the operations related to a user
type User struct{}

// Instantiate a new object of User type
func New() *User {
	return &User{}
}

// topKUsersByPRsAndCommits returns Top K active users sorted
// by amount of PRs created and commits pushed
func (u *User) topKUsersByPRsAndCommits(num int) error {
	return nil
}

func (u *User) CmdTopKUsersByPRsAndCommits() *cli.Command {
	return &cli.Command{
		Name:    "topk-by-pc",
		Aliases: []string{"t"},
		Usage:   "Top K active users sorted by amount of PRs created and commits pushed",
		Flags:   []cli.Flag{},
		Action: func(*cli.Context) error {
			k := 10
			return u.topKUsersByPRsAndCommits(k)
		},
	}
}
