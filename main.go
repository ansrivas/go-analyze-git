package main

import (
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	figure "github.com/common-nighthawk/go-figure"
	config "gitlab.com/ansrivas/go-analyze-git/internal/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// BuildTime gets populated during the build proces
	BuildTime = ""

	//Version gets populated during the build process
	Version = ""
)

// setupConfigOrFatal loads all the variables from the environment variable.
// At this point everything is read as a Key,Value in a map[string]string
func setupConfigOrFatal() config.Config {
	conf, err := config.LoadEnv()
	if err != nil {

		log.Fatal().Msgf("Failed to parse the environment variable. Error %s", err.Error())
	}
	return conf
}

func printBanner() {
	myFigure := figure.NewFigure("go-analyze-git", "", true)
	myFigure.Print()
}

// setupLogger will setup the zap json logging interface
// if the --debug flag is passed, level will be debug
func setupLogger(debug bool) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

}

// printVersionInfo just prints the build time and git commit/tag used
// for this build
func printVersionInfo(version bool) {
	if version {
		fmt.Printf("Version  : %s\nBuildTime: %s\n", Version, BuildTime)
		os.Exit(0)
	}
}

func helloWorld() string {
	return "Hello World"
}

func main() {
	debug := flag.Bool("debug", false, "Set the log level to debug")
	version := flag.Bool("version", false, "Display the BuildTime and Version of this binary")
	flag.Parse()

	printVersionInfo(*version)
	setupLogger(*debug)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println("Press ctrl-c to exit")
	log.Info().Msgf("Exiting server. Message: %v", <-errc)
}
