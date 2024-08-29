/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print spicetify config",
	Run: func(cmd *cobra.Command, args []string) {
		rootLogger.Infof("config file path: %s", paths.ConfigPath)
		rootLogger.Infof("daemon: %t", vars.Daemon)
		rootLogger.Infof("mirror: %t", vars.Mirror)
		rootLogger.Infof("Spotify data path: %s", vars.SpotifyDataPath)
		rootLogger.Infof("Spotify exec path: %s", vars.SpotifyExecPath)
		rootLogger.Infof("Spotify config path: %s", vars.SpotifyConfigPath)
	},
}
