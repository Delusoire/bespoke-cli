/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package vars

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	Mirror            bool
	SpotifyDataPath   string
	SpotifyExecPath   string
	SpotifyConfigPath string
	CfgFile           string
)

var (
	Daemon bool
)

var RootLogger = log.New(os.Stderr)
