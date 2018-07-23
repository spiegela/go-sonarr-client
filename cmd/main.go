package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	homeFolderName = ".go-sonarr-client"
)

func main() {
	app := cli.NewApp()

	app.Name = "sonarr"
	app.Usage = "Use Sonarr from the terminal"
	app.Version = "0.0.1"

	// create a new goroutine to gracefully shutdown datastore
	// shutdown := make(chan int, 1)

	// wg := sync.WaitGroup{}

	// cli.OsExiter = func(c int) {
	// 	shutdown <- c
	// }

	// wg.Add(1)

	// go func(shutdown chan int) {
	// 	defer wg.Done()
	// 	exitCode := <-shutdown

	// 	db.Close()
	// 	os.Exit(exitCode)
	// }(shutdown)

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "save",
			Usage:  "save will prompt you to enter the url and api key for sonarr",
			Action: save,
		},
		cli.Command{
			Name:   "credentials",
			Usage:  "print your sonarr credentials",
			Action: getCredentials,
		},
		cli.Command{
			Name:   "search",
			Usage:  "search tvdb for a show",
			Action: search,
		},
		cli.Command{
			Name:   "show",
			Usage:  "display series info",
			Action: showSeriesInfo,
		},
		cli.Command{
			Name:   "add",
			Usage:  "add a series to wanted",
			Action: addSeries,
		},
		// cli.Command{
		// 	Name:   "delete",
		// 	Usage:  "delete a movie from wanted",
		// 	Action: deleteMovie,
		// },
		cli.Command{
			Name:   "unlock",
			Usage:  "remove the PID lock file in the case a bug interrupts our app",
			Action: unlock,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("cli for sonarr failed to run: %v\n", err)
		os.Exit(1)
	}

	// shutdown <- 0

	// wg.Wait()
}
