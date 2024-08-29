/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spotify

import (
	"github.com/Delusoire/bespoke-cli/v3/cmd/vars"
	"github.com/spf13/cobra"
)

var rootLogger = vars.RootLogger

func AddCommands(c *cobra.Command) {
	c.AddCommand(runCmd)
	c.AddCommand(updateCmd)
}
