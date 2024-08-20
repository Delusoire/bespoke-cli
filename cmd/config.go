/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print spicetify config",
	Run: func(cmd *cobra.Command, args []string) {
		rootLogger.Infof("config file path: %s", paths.ConfigPath)
		rootLogger.Infof("daemon: %s", daemon)
		rootLogger.Infof("mirror: %s", mirror)
		rootLogger.Infof("Spotify data path: %s", spotifyDataPath)
		rootLogger.Infof("Spotify exec path: %s", spotifyExecPath)
		rootLogger.Infof("Spotify config path: %s", spotifyConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
