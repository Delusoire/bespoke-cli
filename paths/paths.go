/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package paths

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	ConfigPath = GetDefaultSpicetifyConfigPath()
)

func GetDefaultSpicetifyConfigPath() string {
	// <portableAppPath>/
	// ├─ bin/
	// │  ├─ <realExe>
	// ├─ config/

	// <portableAppPath> should case-insensitively match "**/spicetify*"
	if exe, err := os.Executable(); err == nil {
		if realExe, err := filepath.EvalSymlinks(exe); err == nil {
			portableBinPath := filepath.Dir(realExe)
			portableAppPath := filepath.Dir(portableBinPath)

			portableBinDir := filepath.Base(portableAppPath)
			portableAppDir := strings.ToLower(filepath.Base(portableAppPath))

			if strings.HasPrefix(portableAppDir, "spicetify") && portableBinDir == "bin" {
				portableConfigPath := filepath.Join(portableAppPath, "config")

				if EnsurePath(portableConfigPath) {
					return portableConfigPath
				}
			}
		}
	}

	return GetPlatformSpicetifyConfigPath()
}

func GetDefaultSpotifyDataPath() string {
	path, err := GetPlatformSpotifyDataPath()
	if err != nil {
		panic(err)
	}
	return path
}

func GetDefaultSpotifyExecPath(spotifyDataPath string) string {
	return GetPlatformSpotifyExecPath(spotifyDataPath)
}

func GetDefaultSpotifyConfigPath() string {
	path, err := GetPlatformSpotifyConfigPath()
	if err != nil {
		panic(err)
	}
	return path
}

func GetSpotifyAppsPath(spotifyPath string) string {
	return filepath.Join(spotifyPath, "Apps")
}

func ResolveHomePaths(relPaths ...string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	paths := make([]string, len(relPaths))
	for i, relPath := range relPaths {
		paths[i] = filepath.Join(home, relPath)
	}
	return paths
}

func EnsurePath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
