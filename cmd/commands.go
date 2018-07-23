package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sonarr "github.com/jrudio/go-sonarr-client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

func startDB() (store, error) {
	// create persistent key store in user home directory
	storeDirectory, err := homedir.Dir()

	if err != nil {
		return store{}, err
	}

	storeDirectory = filepath.Join(storeDirectory, homeFolderName)

	return initDataStore(storeDirectory)
}

func unlock(c *cli.Context) error {
	storeDirectory, err := homedir.Dir()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	storeDirectory = filepath.Join(storeDirectory, homeFolderName)
	lockFilePath := filepath.Join(storeDirectory, "LOCK")

	if err := os.Remove(lockFilePath); err != nil {
		return cli.NewExitError(fmt.Sprintf("failed to remove file: %v", err), 1)
	}

	fmt.Println("removed LOCK file")

	return nil
}

func save(c *cli.Context) error {
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// prompt to save url to sonarr application
	fmt.Println("enter the url that points to sonarr...")

	var sonarrURL string

	fmt.Scanln(&sonarrURL)

	if sonarrURL == "" {
		return cli.NewExitError("url is required", 1)
	}

	// prompt to save api key
	fmt.Println("enter your api key...")

	var key string

	fmt.Scanln(&key)

	if key == "" {
		return cli.NewExitError("api key is required", 1)
	}

	// confirm
	fmt.Printf("URL: %s\nAPI Key: %s\n", sonarrURL, key)
	fmt.Println("Are you sure you want to save?")

	// show success/error

	if err := db.saveSonarrURL(sonarrURL); err != nil {
		return cli.NewExitError(fmt.Sprintf("save url failed: %v", err), 1)
	}

	if err := db.saveSonarrKey(key); err != nil {
		// revert url save
		db.saveSonarrURL("")

		return cli.NewExitError(fmt.Sprintf("save api key failed: %v", err), 1)
	}

	return nil
}

func getCredentials(c *cli.Context) error {
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	radarrURL, err := db.getSonarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	key, err := db.getSonarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("URL: %s\nAPI Key: %s\n", radarrURL, key)

	return nil
}

func search(c *cli.Context) error {
	title := strings.Join(c.Args(), " ")

	if title == "" {
		return cli.NewExitError("a title is required", 1)
	}

	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	sonarrKey, err := db.getSonarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sonarrURL, err := db.getSonarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	client, err := sonarr.New(sonarrURL, sonarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	results, err := client.Search(title)

	for _, series := range results {
		fmt.Printf("%s (%d) - %d\n", series.Title, series.Year, series.TvdbID)
	}

	return nil
}

func showSeriesInfo(c *cli.Context) error {
	tmdbIDstr := c.Args().First()

	if tmdbIDstr == "" {
		return cli.NewExitError("a tvdb id is required", 1)
	}

	// fire up store
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// grab credentials
	sonarrKey, err := db.getSonarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sonarrURL, err := db.getSonarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// create sonarr client to interface with sonarr
	client, err := sonarr.New(sonarrURL, sonarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	tvdbID, err := strconv.Atoi(tmdbIDstr)

	series, err := client.GetSeriesFromTVDB(tvdbID)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// title (year) - tmdbid
	// 		summary
	const output = "%s (%d) - %d\n\t%s\n"

	fmt.Printf(output, series.Title, series.Year, series.TvdbID, series.Overview)

	return nil
}

func addSeries(c *cli.Context) error {
	tvdbIDstr := c.Args().First()

	if tvdbIDstr == "" {
		return cli.NewExitError("a tvdb id is required", 1)
	}

	// fire up store
	db, err := startDB()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	defer db.Close()

	// grab credentials
	sonarrKey, err := db.getSonarrKey()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	sonarrURL, err := db.getSonarrURL()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// create sonarr client to interface with sonarr
	client, err := sonarr.New(sonarrURL, sonarrKey)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	tmdbIDStr, err := strconv.Atoi(tvdbIDstr)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	series, err := client.GetSeriesFromTVDB(tmdbIDStr)

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	// show available profiles
	profiles, err := client.GetProfiles()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	profileCount := len(profiles)

	if profileCount == 0 {
		fmt.Println("aborting...")
		return cli.NewExitError("no profiles found", 1)
	}

	fmt.Print("available quality profiles:\n\n")

	for i, profile := range profiles {
		fmt.Printf("[%d] - %s\n", i, profile.Name)
	}

	fmt.Print("\nplease choose a profile: ")

	// ask user for requested quality
	var requestedQualityIndex int
	fmt.Scanln(&requestedQualityIndex)

	// bound-check user input
	if requestedQualityIndex < 0 || requestedQualityIndex > profileCount {
		return cli.NewExitError("invalid selection", 1)
	}

	profile := profiles[requestedQualityIndex].ID

	// display available root folders
	folders, err := client.GetRootFolders()

	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(folders) == 0 {
		fmt.Println("aborting...")
		return cli.NewExitError("failed to find root folders", 1)
	}

	fmt.Println("\navailable root folders:")

	for i, folder := range folders {
		fmt.Printf("[%d] - %s\n", i, folder.Path)
	}

	fmt.Print("\nchoose a folder to download this series to: ")

	// ask user where we should download this series to
	var rootFolderPathIndex int
	fmt.Scanln(&rootFolderPathIndex)

	fmt.Println()

	rootFolder := folders[rootFolderPathIndex].Path

	// set movie path and profile quality to user preference
	series.AddOptions.SearchForMissingEpisodes = true
	series.QualityProfileID = profile
	series.RootFolderPath = rootFolder
	series.Monitored = true

	if err := client.AddSeries(*series); err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("added %s (%d) successfully\n", series.Title, series.Year)

	return nil
}

func deleteMovie(c *cli.Context) error {
	return nil
}
