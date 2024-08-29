/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/Delusoire/bespoke-cli/v3/paths"
	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix your spotify installation",
	Run: func(cmd *cobra.Command, args []string) {
		if err := execFix(rootLogger); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Restored Spotify to stock state")
	},
}

func execFix(logger *log.Logger) error {
	if vars.Mirror {
		os.RemoveAll(filepath.Join(paths.ConfigPath, "apps"))
	} else {
		spaBakGlob := filepath.Join(paths.GetSpotifyAppsPath(vars.SpotifyDataPath), "*.spa.bak")
		spaBaks, err := filepath.Glob(spaBakGlob)
		if err != nil {
			return err
		}
		if len(spaBaks) == 0 {
			return fmt.Errorf("Spotify is already in stock state!")
		}

		for _, spaBak := range spaBaks {
			spa := strings.TrimSuffix(spaBak, ".bak")
			if err = os.RemoveAll(strings.TrimSuffix(spa, ".spa")); err != nil {
				logger.Warn(err)
			}
			if err = os.Rename(spaBak, spa); err != nil {
				logger.Errorf("failed to restore %s: %s", spaBak, err)
			}
		}
	}

	return nil
}
