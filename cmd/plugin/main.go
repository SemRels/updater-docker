// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	plugin "github.com/SemRels/updater-docker/internal/plugin"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Getenv))
}

func run(stdout, stderr io.Writer, getenv func(string) string) int {
	version := getenv("SEMREL_VERSION")
	if version == "" {
		version = getenv("SEMREL_NEXT_VERSION")
	}
	if version == "" {
		fmt.Fprintln(stderr, "updater-docker: SEMREL_VERSION is required")
		return 1
	}
	version = strings.TrimPrefix(version, "v")

	file := getenv("SEMREL_PLUGIN_FILE")
	if file == "" {
		file = "Dockerfile"
	}
	image := getenv("SEMREL_PLUGIN_IMAGE")
	registry := getenv("SEMREL_PLUGIN_REGISTRY")

	if getenv("SEMREL_DRY_RUN") == "true" {
		fmt.Fprintf(stdout, "updater-docker: [dry-run] would update %s to version %s\n", file, version)
		return 0
	}

	if err := plugin.NewUpdater().Update(file, image, registry, version); err != nil {
		fmt.Fprintln(stderr, "updater-docker:", err)
		return 1
	}

	fmt.Fprintf(stdout, "updater-docker: updated %s to version %s\n", file, version)
	return 0
}
