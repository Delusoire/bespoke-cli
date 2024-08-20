/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"fmt"

	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
)

var (
	showSpotiyDataPath    bool
	showSpotifyExecPath   bool
	showSpotifyConfigPath bool
	showConfigPath        bool
)

var pathsCmd = &cobra.Command{
	Use:   "paths",
	Short: "Print spicetify config",
	Run: func(cmd *cobra.Command, args []string) {
		if !showSpotiyDataPath && !showSpotifyExecPath && !showSpotifyConfigPath && !showConfigPath {
			showSpotiyDataPath = true
			showSpotifyExecPath = true
			showSpotifyConfigPath = true
			showConfigPath = true
		}
		fmt.Println("mirror:", mirror)
		if showSpotiyDataPath {
			fmt.Println("Spotify data path:", spotifyDataPath)
		}
		if showSpotifyExecPath {
			fmt.Println("Spotify exec path:", spotifyExecPath)
		}
		if showSpotifyConfigPath {
			fmt.Println("Spotify config path:", spotifyConfigPath)
		}
		if showConfigPath {
			fmt.Println("config file path:", paths.ConfigPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(pathsCmd)

	pathsCmd.Flags().BoolVar(&showSpotiyDataPath, "spotify-data-path", false, "Show Spotify data path")
	pathsCmd.Flags().BoolVar(&showSpotifyExecPath, "spotify-exec-path", false, "Show Spotify exec path")
	pathsCmd.Flags().BoolVar(&showSpotifyConfigPath, "spotify-config-path", false, "Show Spotify config path")
	pathsCmd.Flags().BoolVar(&showConfigPath, "config", false, "Show config path")
}
