/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spotify

import (
	"os/exec"
	"path/filepath"

	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/Delusoire/bespoke-cli/v3/paths"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch Spotify with your favorite addons",
	Run: func(cmd *cobra.Command, args []string) {
		execRun(args)
	},
}

func prepend[Type any](slice []Type, elems ...Type) []Type {
	return append(elems, slice...)
}

func execRun(args []string) {
	defaultArgs := []string{ /*"--disable-web-security",*/ }
	args = prepend(args, defaultArgs...)
	if vars.Mirror {
		args = prepend(args, "--app-directory="+filepath.Join(paths.ConfigPath, "apps"))
	}
	exec.Command(vars.SpotifyExecPath, args...).Start()
}
