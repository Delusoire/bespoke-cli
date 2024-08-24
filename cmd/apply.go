/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Delusoire/bespoke-cli/v3/archive"
	"github.com/Delusoire/bespoke-cli/v3/link"
	"github.com/Delusoire/bespoke-cli/v3/paths"
	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply spicetify patches on Spotify",
	Run: func(cmd *cobra.Command, args []string) {
		if err := execApply(rootLogger); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Patched Spotify")
	},
}

func getApps() (src string, dest string) {
	src = paths.GetSpotifyAppsPath(spotifyDataPath)
	if mirror {
		dest = filepath.Join(paths.ConfigPath, "apps")
	} else {
		dest = src
	}
	return src, dest
}

func extractSpa(spa string, destFolder string, logger *log.Logger) error {
	basename := filepath.Base(spa)
	extractDest := filepath.Join(destFolder, strings.TrimSuffix(basename, ".spa"))
	logger.Infof("Extracting %s -> %s", spa, extractDest)

	unzipSpa := func(spa, extractDest string) error {
		r, err := zip.OpenReader(spa)
		if err != nil {
			return err
		}
		defer r.Close()

		return archive.UnZip(&r.Reader, extractDest)
	}

	if err := unzipSpa(spa, extractDest); err != nil {
		return err
	}

	if !mirror {
		spaBak := spa + ".bak"
		logger.Infof("Moving %s -> %s", spa, spaBak)

		if err := os.Rename(spa, spaBak); err != nil {
			return err
		}
	}
	return nil
}

func patchFile(path string, patch func(string) string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	content := patch(string(raw))

	return os.WriteFile(path, []byte(content), 0700)
}

func patchIndexHtml(destXpuiPath string, logger *log.Logger) error {
	logger.Info("Patching xpui/index.html")
	return patchFile(filepath.Join(destXpuiPath, "index.html"), func(s string) string {
		return strings.Replace(s, `<script defer="defer" src="/vendor~xpui.js"></script><script defer="defer" src="/xpui.js"></script>`, `<script type="module" src="/hooks/index.js"></script>`, 1)
	})
}

func linkFiles(destXpuiPath string, logger *log.Logger) error {
	folders := []string{"hooks", "modules", "store"}
	for _, folder := range folders {
		folderSrcPath := filepath.Join(paths.ConfigPath, folder)
		folderDestPath := filepath.Join(destXpuiPath, folder)
		logger.Infof("Linking %s -> %s", folderDestPath, folderSrcPath)

		os.Remove(folderDestPath)
		if err := link.Create(folderSrcPath, folderDestPath); err != nil {
			return err
		}
	}
	return nil
}

func execApply(logger *log.Logger) error {
	src, dest := getApps()

	spa := filepath.Join(src, "xpui.spa")
	if err := extractSpa(spa, dest, logger); err != nil {
		return fmt.Errorf("failed to extract xpui.spa: %w", err)
	}

	destXpuiPath := filepath.Join(dest, "xpui")
	if err := patchIndexHtml(destXpuiPath, logger); err != nil {
		return fmt.Errorf("failed to patch index.html: %w", err)
	}

	if err := linkFiles(destXpuiPath, logger); err != nil {
		return fmt.Errorf("failed to link files: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
