//go:build linux

/*
 * Copyright (C) 2024 ririxi, Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package paths

import (
	"path/filepath"

	e "github.com/Delusoire/bespoke-cli/v3/errors"

	"github.com/adrg/xdg"
)

func isSpotifyDataPath(path string) bool {
	return EnsurePath(GetPlatformSpotifyExecPath(path))
}

func GetPlatformSpotifyDataPath() (string, error) {
	absPaths := []string{
		"/opt/spotify/",
		"/opt/spotify/spotify-client/",
		"/usr/share/spotify/",
		"/usr/libexec/spotify/",
		"/var/lib/flatpak/app/com.spotify.Client/x86_64/stable/active/files/extra/share/spotify/",
	}

	homePaths := ResolveHomePaths(
		".local/share/flatpak/app/com.spotify.Client/x86_64/stable/active/files/extra/share/spotify/",
		".local/share/spotify-launcher/install/usr/share/spotify/",
	)

	paths := append(absPaths, homePaths...)

	for _, path := range paths {
		if isSpotifyDataPath(path) {
			return path, nil
		}
	}

	return "", e.ErrPathNotFound
}

func GetPlatformSpotifyExecPath(spotifyDataPath string) string {
	return filepath.Join(spotifyDataPath, "spotify")
}

func GetPlatformSpotifyConfigPath() (string, error) {
	pref := filepath.Join(xdg.CacheHome, "spotify")

	if EnsurePath(pref) {
		return pref, nil
	}

	return "", e.ErrPathNotFound
}

func GetPlatformSpicetifyConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "spicetify")
}
