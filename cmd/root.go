/*
Copyright © 2022 Zack Teska <zerodahero@gmail.com>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//go:generate bash ../version.sh
//go:embed version.txt
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "lester",
	Short:   "The Legacy Tester",
	Version: strings.TrimSpace(version),
	Long:    ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initConfig)

	viper.SetDefault("author", "Zack Teska <zerodahero@gmail.com>")
	viper.SetConfigName("lester")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
}

func initConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
