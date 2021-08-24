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
	"container/heap"
	"os"
	"strings"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/internal/flags"
	"gitlab.com/ansrivas/go-analyze-git/internal/utils"
	"gitlab.com/ansrivas/go-analyze-git/pkg/events"
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
func (u *User) topKUsersByPRsAndCommits(count int, actorsFile, eventsFile, commitsFile string) (utils.GenericDictHeap, error) {

	actorsChan := make(chan string, 10)
	eventsChan := make(chan string, 10)
	commitsChan := make(chan string, 10)
	outputChan := make(chan utils.GenericDictHeap, 1)

	userIDToCommitPRCountsCache := make(map[string]int)
	eventIDToUserIDCache := make(map[string]string)
	userIDToUsernameCache := make(map[string]string)
	var g run.Group

	{
		g.Add(
			utils.ExecuteFunc(actorsFile, actorsChan),
			utils.InterruptFunc("The actorsFile actor was interrupted with: %v\n"))
	}
	{
		g.Add(
			utils.ExecuteFunc(eventsFile, eventsChan),
			utils.InterruptFunc("The eventsFile actor was interrupted with: %v\n"))
	}
	{
		g.Add(
			utils.ExecuteFunc(commitsFile, commitsChan),
			utils.InterruptFunc("The commitsFile actor was interrupted with: %v\n"))
	}

	{
		g.Add(func() error {
			gdHeap := &utils.GenericDictHeap{}
			heap.Init(gdHeap)

			for event := range eventsChan {
				columns := strings.Split(event, ",")
				if len(columns) != 4 {
					continue
				}
				if columns[1] != events.Push && columns[1] != events.Create {
					// In case events type is not "PushEvent" or "CreateEvent"
					continue
				}
				userID := columns[2]
				eventID := columns[0]
				eventIDToUserIDCache[eventID] = userID
				userIDToCommitPRCountsCache[userID] = 0
			}
			for commitLine := range commitsChan {
				columns := strings.Split(commitLine, ",")
				eventID := columns[len(columns)-1]
				// Check if this eventID is present in eventIDToUserIDCache and is a valid PushEvent
				userID, ok := eventIDToUserIDCache[eventID]
				if !ok {
					continue
				}
				userIDToCommitPRCountsCache[userID] += 1
			}

			for line := range actorsChan {
				columns := strings.Split(line, ",")
				userIDToUsernameCache[columns[0]] = columns[1]
			}

			// Now iterate over userIDToCommitPRCountsCache and
			// convert the userID -> userName and populate the heap
			for userID, commitCount := range userIDToCommitPRCountsCache {
				userName, exists := userIDToUsernameCache[userID]
				if !exists {
					continue
				}

				heap.Push(gdHeap, utils.GenericDict{
					Key:   userName,
					Value: commitCount,
				})
				// Maintaining only top-k elements
				if gdHeap.Len() > count {
					heap.Pop(gdHeap)
				}
				delete(userIDToUsernameCache, userID)
			}
			initialHeapLen := gdHeap.Len()
			result := make(utils.GenericDictHeap, initialHeapLen)
			for i := initialHeapLen; i > 0; i-- {
				result[i-1] = heap.Pop(gdHeap).(utils.GenericDict)
			}
			outputChan <- result
			return nil
		}, utils.InterruptFunc("The worker actor was interrupted with: %v\n"))
	}

	err := g.Run()
	return <-outputChan, err
}

func (u *User) CmdTopKUsersByPRsAndCommits() *cli.Command {
	cmdName := "topk-by-pc"
	return &cli.Command{
		Name:    cmdName,
		Aliases: []string{"t"},
		Usage:   "Top K active users sorted by amount of PRs created and commits pushed",
		Flags: []cli.Flag{
			flags.CommitsFileFlag,
			flags.EventsFileFlag,
			flags.ActorsFileFlag,
			flags.CountFlag,
			flags.JsonFlag,
		},
		Action: func(c *cli.Context) error {
			commitsFile := c.String("commits-file")
			eventsFile := c.String("events-file")
			actorsFile := c.String("actors-file")
			count := c.Int("count")
			json := c.Bool("json")
			start := time.Now()
			output, err := u.topKUsersByPRsAndCommits(count, actorsFile, eventsFile, commitsFile)
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
