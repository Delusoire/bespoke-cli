/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/spf13/cobra"
)

var rootLogger = vars.RootLogger

func AddCommands(c *cobra.Command) {
	c.AddCommand(applyCmd)
	c.AddCommand(configCmd)
	c.AddCommand(daemonCmd)
	c.AddCommand(devCmd)
	c.AddCommand(fixCmd)
	c.AddCommand(initCmd)
	c.AddCommand(pkgCmd)
	c.AddCommand(protocolCmd)
	c.AddCommand(syncCmd)
}
