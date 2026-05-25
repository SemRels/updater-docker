// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

package main

import (
	"log"

	plugin "github.com/SemRels/updater-docker/internal/plugin"
)

func main() {
	tagger := plugin.NewTagger(plugin.Config{})
	log.Printf("updater-docker plugin ready: tags and pushes Docker images (%T)", tagger)
}
