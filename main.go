package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/ansrivas/go-analyze-git/pkg/app"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// BuildTime gets populated during the build proces
	BuildTime = ""

	//Version gets populated during the build process
	Version = ""
)

// setupLogger will setup the zap json logging interface
// if the --debug flag is passed, level will be debug
func setUpLogger(c *cli.Context) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if c.Bool("debug") {
		fmt.Fprintf(c.App.Writer, "Setting the log level to debug")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func createApp() *app.App {
	cliApp := app.New()
	cliApp.Name = "go-analyze-git"
	cliApp.Description = "Run analytics on git data set"
	cliApp.Version = Version

	cliApp.Authors = []*cli.Author{
		{
			Name:  "Ankur Srivastava",
			Email: "ankur.srivastava@email.de",
		},
	}
	cliApp.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Set the log level to debug",
		},
		&cli.BoolFlag{
			Name:  "buildtime",
			Usage: "time of this build",
		},
	}
	cliApp.Before = func(c *cli.Context) error {
		if c.IsSet("buildtime") {
			fmt.Fprintf(c.App.Writer, "%v buildtime: %v\n", cliApp.Name, BuildTime)
			return cli.Exit("", 0)
		}
		// Set up log level
		setUpLogger(c)
		return nil
	}
	cliApp.EnableBashCompletion = true
	cliApp.Commands = []*cli.Command{
		cliApp.User(),
		cliApp.Repository(),
	}
	sort.Sort(cli.CommandsByName(cliApp.Commands))
	return cliApp

}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := createApp()
	err := app.RunWithContext(ctx, os.Args)
	if err != nil {
		log.Error().Msgf("Error in app run: %s", err)
	}

	// errc := make(chan error)
	// go func() {
	// 	c := make(chan os.Signal, 2)
	// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// 	errc <- fmt.Errorf("%s", <-c)
	// }()
	// fmt.Println("Press ctrl-c to exit")
	// log.Info().Msgf("Exiting server. Message: %v", <-errc)
}
