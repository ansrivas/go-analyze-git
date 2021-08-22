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

package app

import (
	"context"

	"github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/pkg/repository"
	"gitlab.com/ansrivas/go-analyze-git/pkg/user"
)

type App struct {
	*cli.App
}

func New() *App {
	return &App{
		cli.NewApp(),
	}
}
func (c *App) Repository() *cli.Command {
	repo := repository.New()
	return &cli.Command{
		Name:    "repository",
		Aliases: []string{"r"},
		Usage:   "Commands related to repository operations",
		Action: func(*cli.Context) error {
			return nil
		},
		Subcommands: []*cli.Command{
			repo.CmdTopKReposByCommits(),
			repo.CmdTopKReposByWatchEvents(),
		},
	}
}

func (c *App) User() *cli.Command {
	userCli := user.New()
	return &cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "Commands related to user operations",
		Action: func(*cli.Context) error {
			return nil
		},
		Subcommands: []*cli.Command{
			userCli.CmdTopKUsersByPRsAndCommits(),
		},
	}
}

// RunWithContext is a wrapper on urfave/cli RunContext function
func (a *App) RunWithContext(ctx context.Context, arguments []string) error {
	return a.RunContext(ctx, arguments)
}
