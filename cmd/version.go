package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version    = "unknown"
	commit     = "none"
	buildDate  = "unknown"
	builtBy    = "unknown"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of PoDownloader",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("PoDownloader version: %s, commit:%s", version, commit))
			fmt.Println(fmt.Sprintf("Built at %s by %s", buildDate, builtBy))
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
