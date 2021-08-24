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

// topKReposByCommits returns Top K repositories by
// the amount of commits pushed
func (r *Repository) topKReposByCommits(count int, reposFile, eventsFile, commitsFile string) (utils.GenericDictHeap, error) {

	reposChan := make(chan string, 10)
	eventsChan := make(chan string, 10)
	commitsChan := make(chan string, 10)
	outputChan := make(chan utils.GenericDictHeap, 1)
	var g run.Group

	{
		g.Add(
			executeFunc(reposFile, reposChan),
			interruptFunc("The reposFile actor was interrupted with: %v\n"))
	}

	{
		g.Add(
			executeFunc(eventsFile, eventsChan),
			interruptFunc("The eventsFile actor was interrupted with: %v\n"))
	}

	{
		g.Add(
			executeFunc(commitsFile, commitsChan),
			interruptFunc("The commitsFile actor was interrupted with: %v\n"))
	}

	{
		g.Add(func() error {
			eventsToRepoCache := make(map[string]string)

			// Valid repos
			repoToCommitsCountCache := make(map[string]int)

			for line := range eventsChan {
				columns := strings.Split(line, ",")
				if len(columns) != 4 {
					log.Debug().Msgf("Missing 4 columns in [%s]", line)
					continue
				}

				eventType := columns[1]
				// Filter out all the PushEvents
				if eventType != events.Push {
					continue
				}
				repoID := columns[3]
				eventID := columns[0]
				eventsToRepoCache[eventID] = repoID
				repoToCommitsCountCache[repoID] = 0
			}

			for commitLine := range commitsChan {
				// Although commits are supposed to be 3 entries, but
				// there is a chance they might not be.
				// We will make an assumption that in a line, we test
				// if the entry is an integer - we use it.
				// TODO (ansrivas): Benchmark with csv-scan split function as well
				columns := strings.Split(commitLine, ",")
				if len(columns) == 0 {
					continue
				}

				commitID := columns[len(columns)-1]
				// Check if this event_id is present in eventsToRepoCache and is a valid PushEvent
				repoID, ok := eventsToRepoCache[commitID]
				if !ok {
					continue
				}
				repoToCommitsCountCache[repoID] += 1
			}

			gdHeap := &utils.GenericDictHeap{}
			heap.Init(gdHeap)
			// For each
			for repoID, commitCount := range repoToCommitsCountCache {
				heap.Push(gdHeap, utils.GenericDict{
					Key:   repoID,
					Value: commitCount,
				})
				// Maintaining only top-k elements
				if gdHeap.Len() > count {
					heap.Pop(gdHeap)
				}
			}

			repoIDToNameCache := make(map[string]string)
			for line := range reposChan {
				columns := strings.Split(line, ",")
				repoIDToNameCache[columns[0]] = columns[1]
			}

			initialHeapLen := gdHeap.Len()
			result := make(utils.GenericDictHeap, initialHeapLen)
			for i := initialHeapLen; i > 0; i-- {
				gd := heap.Pop(gdHeap).(utils.GenericDict)
				repoName, ok := repoIDToNameCache[gd.Key]
				if !ok {
					log.Error().Msgf("Couldn't find the reponame in cache. Keeping ID")
				} else {
					// mutate gd.Key to store the name, rather than id now
					// Hacky but meh
					gd.Key = repoName
				}
				result[i-1] = gd
			}

			outputChan <- result
			return nil
		}, interruptFunc("The worker actor was interrupted with: %v\n"))
	}

	err := g.Run()
	return <-outputChan, err
}

func (r *Repository) CmdTopKReposByCommits() *cli.Command {

	cmdName := "topk-by-commits"
	return &cli.Command{
		Name:    cmdName,
		Aliases: []string{"tc"},
		Usage:   "Top K repositories by the amount of commits pushed",
		Flags: []cli.Flag{
			flags.ReposFileFlag,
			flags.EventsFileFlag,
			flags.CommitsFileFlag,
			flags.CountFlag,
			flags.JsonFlag,
		},
		Action: func(c *cli.Context) error {
			reposFile := c.String("repos-file")
			eventsFile := c.String("events-file")
			commitsFile := c.String("commits-file")
			count := c.Int("count")
			json := c.Bool("json")
			start := time.Now()
			output, err := r.topKReposByCommits(count, reposFile, eventsFile, commitsFile)
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
