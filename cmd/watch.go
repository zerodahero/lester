/*
Copyright © 2022 Zack Teska <zerodahero@gmail.com>
*/
package cmd

import (
	"fmt"
	"lester/config"
	"lester/files"
	"lester/helpers"
	"lester/observer"
	"lester/tester"
	"math"
	"time"

	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch <component>",
	Short: "Watch for changes and re-run tests",
	Long: `Watch files for the given component and re-run
tests anytime there's a change. Feel free to
pass along your own PHPUnit args as needed.

Examples:

# Watch web component and run full test suite
lester watch web

# Watch root and run only group "stuff" (note the double dash '--' after the component)
lester watch root -- --group banana

# Watch just a single file (relative to root dir)
# This runs for any change in the component
lester watch core -- tests/some/kind/of/test/goes/HereTest.php

# Run all tests matching the filter
lester watch web -- --filter stuff
`,
	Run: func(cmd *cobra.Command, args []string) {
		var component string
		if len(args) == 0 {
			component = "root"
		} else {
			component = args[0]
			args = args[1:]
		}

		componentPath := config.GetComponentPathOrDefault(component)

		fmt.Printf("Watching for files at %s ...\n", componentPath)

		modified := make(chan string)

		go observer.WatchForChanges(componentPath, modified)

		runner := tester.NewTestRunnerFromPath(files.NewFromPath(componentPath))

		// Run initial test
		runner.RunTests(args)

		waitFor := 800 * time.Millisecond
		t := time.AfterFunc(math.MaxInt64, func() {
			helpers.ClearScreen()
			runner.RunTests(args)
		})
		for range modified {
			t.Reset(waitFor)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
