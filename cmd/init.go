/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Delusoire/bespoke-cli/v3/module"
	"github.com/Delusoire/bespoke-cli/v3/paths"
	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Perform one-time spicetify initization",
	Long:  "required to be ran at least once per installation",
	Run: func(cmd *cobra.Command, args []string) {
		if err := execInit(rootLogger); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Initialized spicetify")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func execInit(logger *log.Logger) error {
	configFile := viper.GetViper().ConfigFileUsed()
	if err := viper.SafeWriteConfigAs(configFile); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("failed to write config file: %w", err)
		}
	}

	folders := []string{"hooks", "modules", "store"}
	for _, folder := range folders {
		folderPath := filepath.Join(paths.ConfigPath, folder)
		logger.Debug("Removing folder", "folder", folderPath)
		os.Remove(folderPath)
	}

	if err := module.SetVault(&module.Vault{Modules: map[module.ModuleIdentifier]module.Module{}}); err != nil {
		return fmt.Errorf("failed to initialize vault: %w", err)
	}

	return nil
}
