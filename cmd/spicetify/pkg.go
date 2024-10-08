/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"github.com/Delusoire/bespoke-cli/v3/module"

	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:   "pkg",
	Short: "Manage modules",
}

var pkgInstallCmd = &cobra.Command{
	Use:   "install id url",
	Short: "Add and Install module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		url := args[1]
		identifier := module.NewStoreIdentifier(id)
		aurl := module.ArtifactURL(url).Parse().ToUrl()
		if err := addAndInstall(aurl, identifier); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Module added")
	},
}

func addAndInstall(aurl module.ArtifactURL, identifier module.StoreIdentifier) error {
	if err := module.AddStoreInVault(identifier, &module.Store{
		Installed: false,
		Artifacts: []module.ArtifactURL{aurl},
	}); err != nil {
		return err
	}

	return module.InstallModule(identifier)
}

var pkgDeleteCmd = &cobra.Command{
	Use:   "delete id",
	Short: "Delete and Remove module",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := module.NewStoreIdentifier(args[0])
		if err := deleteAndRemove(identifier); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Module deleted")
	},
}

func deleteAndRemove(identifier module.StoreIdentifier) error {
	if err := module.DeleteModule(identifier); err != nil {
		return err
	}

	return module.RemoveStoreInVault(identifier)
}

var pkgEnableCmd = &cobra.Command{
	Use:   "enable id",
	Short: "Enable or Disable module",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := module.NewStoreIdentifier(args[0])
		if err := module.EnableModuleInVault(identifier); err != nil {
			rootLogger.Fatal(err)
		}
		rootLogger.Info("Module enabled")
	},
}

func init() {
	pkgCmd.AddCommand(pkgInstallCmd, pkgDeleteCmd, pkgEnableCmd)
}
