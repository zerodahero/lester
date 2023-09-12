/*
Copyright © 2022 Zack Teska <zerodahero@gmail.com>
*/
package cmd

import (
	"lester/aggressor"
	"lester/config"
	"lester/files"
	"lester/helpers"
	"lester/observer"
	"lester/printer"
	"math"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var seedFromGit bool = false

// autoCmd represents the auto command
var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automagical Test Watcher",
	Long: `Watches ALL files for changes and will attempt to test all code related
to the modified files for as long as you have this command running.
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := printer.NewPrinter(os.Stdout)
		p.PrintLesterIntro()

		modified := make(chan string)

		config := config.GetConfig()

		a := aggressor.NewAggressor()
		a.InitComponentsFromConfig(config.TestConfig.Components)

		if seedFromGit {
			spinner, _ := pterm.DefaultSpinner.Start("Seeding modified files from git...")
			err := aggressor.SeedFilesFromGit(a, spinner)
			if err != nil {
				pterm.Error.Println(err)
				pterm.Warning.Println("Incomplete seeding from modified files, continuing anyway")
			}
			spinner.Info("Modified files seeded from git")
		}

		pterm.Info.Println("Searching for files to watch automagically...")

		go observer.WatchForChanges(".", modified)

		pterm.Info.Println("Watching for changes...")

		run := make(chan int)
		sem := make(chan int, 1)

		waitFor := 800 * time.Millisecond
		t := time.AfterFunc(math.MaxInt64, func() {
			run <- 1
		})
		for {
			select {
			case event := <-modified:
				a.AddModifiedFile(files.NewFromPath(event))
				t.Reset(waitFor)
			case <-run:
				sem <- 1
				runTests(a)
				<-sem
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(autoCmd)
	autoCmd.Flags().BoolVarP(&seedFromGit, "seed-from-git", "g", false, "Seed initial files from git diff")
}

func runTests(a *aggressor.Aggressor) {
	if !a.HasTestsToRun() {
		pterm.Info.Println("No tests found for modified files :-( ")
		return
	}

	helpers.ClearScreen()
	pterm.Info.Println("Detected modified files:")
	items := []pterm.BulletListItem{}
	for _, file := range a.GetModified() {
		items = append(items, pterm.BulletListItem{Level: 0, Text: file})
	}
	pterm.DefaultBulletList.WithItems(items).Render()

	pterm.DefaultSection.Println("Running tests...")
	if a.RunTests() {
		pterm.Success.Println("Tests complete, nice work! Back to watching...")
	} else {
		pterm.Info.Println("Tests complete.")
	}
}
