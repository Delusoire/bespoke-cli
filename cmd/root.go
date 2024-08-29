/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Delusoire/bespoke-cli/v3/cmd/spicetify"
	"github.com/Delusoire/bespoke-cli/v3/cmd/spotify"
	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "spicetify",
	Short: "Make Spotify your own",
	Long:  `Bespoke is a CLI utility that empowers the desktop Spotify client with custom themes and extensions`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	cobra.MousetrapHelpText = ""
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func init() {
	defaultcfgFile := filepath.Join(paths.ConfigPath, "config.yaml")
	rootCmd.PersistentFlags().StringVar(&vars.CfgFile, "config", defaultcfgFile, "config file")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&vars.Mirror, "mirror", "m", false, "Mirror Spotify files instead of patching them directly")

	defaultSpotifyDataPath := paths.GetDefaultSpotifyDataPath()
	rootCmd.PersistentFlags().StringVar(&vars.SpotifyDataPath, "spotify-data-path", defaultSpotifyDataPath, "Override Spotify data folder")

	defaultSpotifyExecPath := paths.GetDefaultSpotifyExecPath(vars.SpotifyDataPath)
	rootCmd.PersistentFlags().StringVar(&vars.SpotifyExecPath, "spotify-exec-path", defaultSpotifyExecPath, "Override Spotify executable path")

	defaultSpotifyConfigPath := paths.GetDefaultSpotifyConfigPath()
	rootCmd.PersistentFlags().StringVar(&vars.SpotifyConfigPath, "spotify-config-path", defaultSpotifyConfigPath, "Override Spotify config folder (containing prefs & offline.bnk)")

	initViper()

	invokedExecutableName := strings.ToLower(getInvokedExecutableName())
	if strings.HasPrefix(invokedExecutableName, "spotify") {
		spotify.AddCommands(rootCmd)
	} else {
		spicetify.AddCommands(rootCmd)
	}
}

func initViper() {
	viper.BindPFlag("mirror", rootCmd.PersistentFlags().Lookup("mirror"))
	viper.BindPFlag("spotify-data-path", rootCmd.PersistentFlags().Lookup("spotify-data-path"))
	viper.BindPFlag("spotify-exec-path", rootCmd.PersistentFlags().Lookup("spotify-exec-path"))
	viper.BindPFlag("spotify-config-path", rootCmd.PersistentFlags().Lookup("spotify-config-path"))
}

func initConfig() {
	viper.SetConfigFile(vars.CfgFile)
	viper.AutomaticEnv()

	viper.SetDefault("daemon", true)
	viper.SetDefault("mirror", vars.Mirror)
	viper.SetDefault("spotify-data-path", vars.SpotifyDataPath)
	viper.SetDefault("spotify-exec-path", vars.SpotifyExecPath)
	viper.SetDefault("spotify-config-path", vars.SpotifyConfigPath)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		vars.Daemon = viper.GetBool("daemon")
		vars.Mirror = viper.GetBool("mirror")
		vars.SpotifyDataPath = viper.GetString("spotify-data-path")
		vars.SpotifyExecPath = viper.GetString("spotify-exec-path")
		vars.SpotifyConfigPath = viper.GetString("spotify-config-path")
	}
}

func getInvokedExecutableName() string {
	invokedPath := os.Args[0]
	return filepath.Base(invokedPath)
}
