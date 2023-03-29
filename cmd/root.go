/*
Copyright © 2023 Guionardo Furlan <guionardo@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/guionardo/dev-clean/cleaners"
	"github.com/spf13/cobra"
	"github.com/xhit/go-str2duration/v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dev-clean",
	Version: "0.0.1",
	Short:   "Clean projects in [sub]folders",
	Long: `Run this tool into project folders to clean build and
temporary files. For example:

dev-clean /my-project/folder
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return ValidateArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cleaners.AllCleaners().RunClean(cleaner, folderName, minimumAge)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
	// Args: cobra.MinimumNArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	cleanerType   string
	folderName    string
	minimumAgeStr string
	minimumAge    time.Duration
	cleaner       cleaners.Cleaner
)

func init() {
	rootCmd.Flags().StringVarP(&cleanerType, "type", "t", "auto", strings.Join(cleaners.AllCleaners().GetNames(), ", "))
	rootCmd.Flags().StringVarP(&minimumAgeStr, "minimum-age", "a", "7d", "Clean only if last change greater than")
}

func ValidateArgs(args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("Missing target folder argument")
	}
	folderName = args[0]
	if stat, err := os.Stat(folderName); err != nil || !stat.IsDir() {
		return fmt.Errorf("Invalid or not found folder %s", folderName)
	}
	if cleaner, err = cleaners.AllCleaners().GetCleanerFromName(cleanerType, folderName); err != nil {
		return
	}
	if minimumAge, err = str2duration.ParseDuration(minimumAgeStr); err != nil {
		return
	}
	return
}
