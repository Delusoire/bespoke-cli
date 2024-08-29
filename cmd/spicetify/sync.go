/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Delusoire/bespoke-cli/v3/archive"
	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update spicetify hooks from GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		if err := installHooks(); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Hooks updated successfully")
	},
}

// TODO: let the user choose which release to install (& include version compatibility info)
func installHooks() error {
	res, err := http.Get("http://github.com/spicetify/hooks/releases/latest/download/hooks.tar.gz")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	hooksPath := filepath.Join(paths.ConfigPath, "hooks")
	if err := os.RemoveAll(hooksPath); err != nil {
		return err
	}
	return archive.UnTarGZ(res.Body, hooksPath)
}
