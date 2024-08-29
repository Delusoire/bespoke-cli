/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package spicetify

import (
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Delusoire/bespoke-cli/v3/module"

	e "github.com/Delusoire/bespoke-cli/v3/errors"

	"github.com/spf13/cobra"
)

var protocolCmd = &cobra.Command{
	Use:   "protocol [uri]",
	Short: "Internal protocol handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := HandleProtocol(args[0])
		if res != "" {
			open("spotify:app:rpc:" + res)
		}
		if err != nil {
			rootLogger.Fatal(err)
		}
	},
}

func HandleProtocol(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "spicetify" {
		return "", err
	}
	uuid, action, _ := strings.Cut(u.Opaque, ":")
	response := u.Scheme + ":" + uuid + ":"
	arguments := u.Query()
	err = hp(action, arguments)
	if err == nil {
		response += "1"
	} else {
		response += "0"
	}
	if uuid == "0" {
		response = ""
	}
	return response, err
}

func hp(action string, arguments url.Values) error {
	switch action {
	case "add":
	case "fast-install":
	case "fast-enable":
		_artifacts := arguments["artifacts"]

		identifier := module.NewStoreIdentifier(arguments.Get("id"))
		artifacts := make([]module.ArtifactURL, len(_artifacts))
		for i, a := range _artifacts {
			artifacts[i] = module.ArtifactURL(a).Parse().ToUrl()
		}
		checksum := arguments.Get("checksum")

		if err := module.AddStoreInVault(identifier, &module.Store{
			Installed: false,
			Artifacts: artifacts,
			Checksum:  checksum,
		}); err != nil {
			return err
		}

		if action == "add" {
			return nil
		}

		if err := module.InstallModule(identifier); err != nil {
			return err
		}

		if action == "fast-install" {
			return nil
		}

		return module.EnableModuleInVault(identifier)

	case "install":
		identifier := module.NewStoreIdentifier(arguments.Get("id"))
		return module.InstallModule(identifier)

	case "enable":
		identifier := module.NewStoreIdentifier(arguments.Get("id"))
		return module.EnableModuleInVault(identifier)

	case "delete":
		identifier := module.NewStoreIdentifier(arguments.Get("id"))
		return module.DeleteModule(identifier)

	case "remove":
		identifier := module.NewStoreIdentifier(arguments.Get("id"))
		return module.RemoveStoreInVault(identifier)

	case "fast-delete":
	case "fast-remove":
		identifier := module.NewStoreIdentifier(arguments.Get("id"))

		if err := module.EnableModuleInVault(module.StoreIdentifier{
			ModuleIdentifier: identifier.ModuleIdentifier,
		}); err != nil {
			return err
		}

		if err := module.DeleteModule(identifier); err != nil {
			return err
		}

		if action == "fast-delete" {
			return nil
		}

		return module.RemoveStoreInVault(identifier)

	}
	return e.ErrUnsupportedOperation
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
