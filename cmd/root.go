package main

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "PoDownloader is a simple CLI tool to download podcast.",
	Long: `PoDownloader is a simple CLI tool to download podcast.

This tool will download podcast RSS, podcast cover image, episode audio files, episode cover images and episode shownotes.

Use the HTTP_PROXY environment variable to set a HTTP or SOCSK5 proxy.`,
}
