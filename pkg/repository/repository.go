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
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/internal/fileops"
	"gitlab.com/ansrivas/go-analyze-git/internal/flags"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
	"gitlab.com/ansrivas/go-analyze-git/pkg/events"
)

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

// topKReposByEvents returns Top K repositories
// sorted by watch events
func (r *Repository) topKReposByEvents(count int, event, eventsFile, reposFile string) (utils.GenericIntDictList, error) {

	// Read repos file
	reposChan := make(chan string, 10)
	errReposChan := make(chan error)
	fops := fileops.NewWithBufSize(128 * 1024)
	go fops.ReadFileStreaming(reposFile, reposChan, errReposChan, bufio.ScanLines)

	// Read events file
	eventsChan := make(chan string, 10)
	errEventsChan := make(chan error)
	go fops.ReadFileStreaming(eventsFile, eventsChan, errEventsChan, bufio.ScanLines)

	watchEventsCache := make(map[string]int)
	for line := range eventsChan {
		columns := strings.Split(line, ",")
		if len(columns) != 4 {
			log.Debug().Msgf("Missing 4 columns in [%s]", line)
			continue
		}

		eventType := columns[1]
		eventID := columns[3]
		// Filter out all the WatchEvents
		if eventType != events.Watch {
			continue
		}
		watchEventsCache[eventID] += 1
	}

	// Now iterate over reposChan to filterout the names
	// from watchEventsCache
	genericDict := utils.GenericIntDictList{}
	for line := range reposChan {
		columns := strings.Split(line, ",")
		if len(columns) != 2 {
			log.Debug().Msgf("ReposChan: Missing 2 columns in [%s]", line)
			continue
		}
		repoName := columns[1]
		repoID := columns[0]

		watchEventCount, exists := watchEventsCache[repoID]
		if exists {
			genericDict = append(genericDict, utils.GenericIntDict{
				Key:   repoName,
				Value: watchEventCount,
			})
			delete(watchEventsCache, repoID)
		}
	}
	// Now sort and take count elements
	genericDict.SortByValue()
	return genericDict.Take(count), nil
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
		Name:    "topk-by-events",
		Aliases: []string{"tw"},
		Usage:   "Top K repositories sorted by events",
		Flags: []cli.Flag{
			flags.ReposFileFlag,
			flags.EventsFileFlag,
			flags.CountFlag,
			flags.EventTypeFlag,
			flags.JsonFlag,
		},
		Action: func(c *cli.Context) error {
			reposFile := c.String("repos-file")
			eventsFile := c.String("events-file")
			eventType := c.String("event-type")
			count := c.Int("count")
			json := c.Bool("json")
			output, err := r.topKReposByEvents(count, eventType, eventsFile, reposFile)
			if err != nil {
				return err
			}

			// If json, print and return
			if json {
				return output.ToJson(os.Stdout)
			}

			utils.RenderTable(output, []string{"RepoID", "Count"})
			return nil
		},
	}
}
