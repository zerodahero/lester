/*
Copyright © 2022 Zack Teska <zerodahero@gmail.com>
*/
package cmd

import (
	"fmt"
	"lester/config"
	"lester/files"
	"lester/tester"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test <component>",
	Short: "Run PHPUnit tests for a given component",
	Long: `Run PHPUnit tests for any component from the comfort
of right here. Pass along your own flags to phpunit
or just enjoy the ride of a full test run.

Examples:

# Full test run for web component
lester test web

# Only the 'banana' group for core (note the double dash '--' after the component)
lester test core -- --group banana

# Test just a single file (relative to root dir)
lester test core -- tests/some/kind/of/test/goes/HereTest.php

# Run all tests matching the filter
lester test web -- --filter stuff
`,
	Run: func(cmd *cobra.Command, args []string) {
		var component string
		if len(args) == 0 {
			component = "root"
		} else {
			component = args[0]
			args = args[1:]
		}

		fmt.Printf("Running tests for %s...\n", component)

		componentPath := config.GetComponentPathOrDefault(component)

		runner := tester.NewTestRunnerFromPath(files.NewFromPath(componentPath))

		runner.RunTests(args)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
