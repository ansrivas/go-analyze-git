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
	"container/heap"
	"os"
	"strings"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/internal/flags"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
	"gitlab.com/ansrivas/go-analyze-git/pkg/events"
)

// topKReposByEvents returns Top K repositories
// sorted by watch events
func (r *Repository) topKReposByEvents(count int, event, eventsFile, reposFile string) (utils.GenericDictHeap, error) {

	reposChan := make(chan string, 10)
	eventsChan := make(chan string, 10)
	outputChan := make(chan utils.GenericDictHeap, 1)

	var g run.Group
	{
		g.Add(func() error {
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
			gdHeap := &utils.GenericDictHeap{}
			heap.Init(gdHeap)

			for line := range reposChan {
				columns := strings.Split(line, ",")
				if len(columns) != 2 {
					log.Debug().Msgf("ReposChan: Missing 2 columns in [%s]", line)
					continue
				}
				repoName := columns[1]
				repoID := columns[0]

				watchEventCount, exists := watchEventsCache[repoID]
				if !exists {
					continue
				}

				heap.Push(gdHeap, utils.GenericDict{
					Key:   repoName,
					Value: watchEventCount,
				})
				// Maintaining only top-k elements
				if gdHeap.Len() > count {
					heap.Pop(gdHeap)
				}
				delete(watchEventsCache, repoID)

			}

			initialHeapLen := gdHeap.Len()
			result := make(utils.GenericDictHeap, initialHeapLen)
			for i := initialHeapLen; i > 0; i-- {
				result[i-1] = heap.Pop(gdHeap).(utils.GenericDict)
			}
			outputChan <- result
			return nil

		}, utils.InterruptFunc("The final goroutine actor was interrupted with: %v\n"))
	}

	{
		g.Add(
			utils.ExecuteFunc(reposFile, reposChan),
			utils.InterruptFunc("The reposFile actor was interrupted with: %v\n"))
	}

	{
		g.Add(
			utils.ExecuteFunc(eventsFile, eventsChan),
			utils.InterruptFunc("The eventsFile actor was interrupted with: %v\n"))
	}

	err := g.Run()
	return <-outputChan, err
}

func (r *Repository) CmdTopKReposByWatchEvents() *cli.Command {
	cmdName := "topk-by-events"
	return &cli.Command{
		Name:    cmdName,
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
			start := time.Now()
			output, err := r.topKReposByEvents(count, eventType, eventsFile, reposFile)
			if err != nil {
				return err
			}

			// If json, print and return
			if json {
				output.ToJson(os.Stdout) //nolint
			} else {
				utils.RenderTable(output, []string{"RepoID", "Count"})
			}

			log.Debug().Msgf("[%s] took %v", cmdName, time.Since(start))
			return nil
		},
	}
}
