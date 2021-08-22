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

package flags

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/pkg/events"
)

var (
	ReposFileFlag = &cli.StringFlag{
		Name:     "repos-file",
		Usage:    "Path to the repos.csv file",
		EnvVars:  []string{"REPOS_FILE"},
		Required: true,
	}
	EventsFileFlag = &cli.StringFlag{
		Name:     "events-file",
		Usage:    "Path to the events.csv file",
		EnvVars:  []string{"EVENTS_FILE"},
		Required: true,
	}
	CountFlag = &cli.IntFlag{
		Name:    "count",
		Usage:   "Maximum count to show",
		Value:   10,
		EnvVars: []string{"COUNT"},
	}
	EventTypeFlag = &cli.StringFlag{
		Name:    "event-type",
		Usage:   "Event type to analyze from ['WatchEvent']",
		Value:   events.Watch,
		EnvVars: []string{"EVENT_TYPE"},
	}
	JsonFlag = &cli.BoolFlag{
		Name:    "json",
		Usage:   "Render the result as json",
		Value:   false,
		EnvVars: []string{"JSON"},
	}
)
