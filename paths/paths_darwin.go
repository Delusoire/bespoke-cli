//go:build darwin

/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package paths

import (
	"path/filepath"

	e "github.com/Delusoire/bespoke-cli/v3/errors"

	"github.com/adrg/xdg"
)

func GetPlatformSpotifyDataPath() (string, error) {
	spotifyAppResourcesPath := "Spotify.app/Contents/Resources"
	spotifyAppLocations := append([]string{"/Applications"}, ResolveHomePaths("Applications")...)

	for _, spotifyAppLocation := range spotifyAppLocations {
		spotifyDataPath := filepath.Join(spotifyAppLocation, spotifyAppResourcesPath)
		if EnsurePath(spotifyDataPath) {
			return spotifyDataPath, nil
		}
	}

	return "", e.ErrPathNotFound
}

func GetPlatformSpotifyExecPath(spotifyPath string) string {
	return filepath.Join(spotifyPath, "Spotify")
}

func GetPlatformSpotifyConfigPath() (string, error) {
	spotifyConfigPath := filepath.Join(xdg.ConfigHome, "Spotify")

	if EnsurePath(spotifyConfigPath) {
		return spotifyConfigPath, nil
	}

	return "", e.ErrPathNotFound
}

func GetPlatformSpicetifyConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "Spicetify")
}
