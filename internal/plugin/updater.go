// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The semrel Authors

// Package plugin updates Dockerfile base image tags in-place.
package plugin

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var fromPattern = regexp.MustCompile(`^(\s*FROM\s+)([^\s]+)(.*)$`)

// Updater updates Dockerfile FROM tags.
type Updater struct{}

// NewUpdater creates an updater.
func NewUpdater() *Updater {
	return &Updater{}
}

// Update rewrites the matching FROM image tag.
func (u *Updater) Update(path, image, registry, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	updated, changed, err := updateContent(string(data), image, registry, version)
	if err != nil {
		return err
	}
	if !changed {
		return fmt.Errorf("no matching FROM image found in %s", path)
	}

	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func updateContent(content, image, registry, version string) (string, bool, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lines := make([]string, 0)
	changed := false

	for scanner.Scan() {
		line := scanner.Text()
		match := fromPattern.FindStringSubmatch(line)
		if match != nil {
			current := match[2]
			if shouldUpdate(current, image) {
				line = match[1] + buildReference(current, image, registry, version) + match[3]
				changed = true
			}
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", false, fmt.Errorf("scan Dockerfile: %w", err)
	}
	return strings.Join(lines, "\n"), changed, nil
}

func shouldUpdate(current, image string) bool {
	if image == "" {
		return true
	}
	return trimRegistryAndTag(current) == trimRegistryAndTag(image)
}

func buildReference(current, image, registry, version string) string {
	name := trimRegistryAndTag(current)
	if image != "" {
		name = trimRegistryAndTag(image)
	}
	if registry != "" {
		return strings.TrimRight(registry, "/") + "/" + name + ":" + version
	}
	prefix := current[:len(current)-len(trimTag(current))]
	if strings.Contains(prefix, "/") && strings.Contains(trimRegistryAndTag(current), "/") {
		return prefix + version
	}
	return name + ":" + version
}

func trimRegistryAndTag(ref string) string {
	ref = trimDigest(ref)
	if idx := strings.LastIndex(ref, ":"); idx != -1 && !strings.Contains(ref[idx+1:], "/") {
		ref = ref[:idx]
	}
	parts := strings.Split(ref, "/")
	if len(parts) > 1 && (strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") || parts[0] == "localhost") {
		return strings.Join(parts[1:], "/")
	}
	return ref
}

func trimTag(ref string) string {
	ref = trimDigest(ref)
	if idx := strings.LastIndex(ref, ":"); idx != -1 && !strings.Contains(ref[idx+1:], "/") {
		return ref[idx+1:]
	}
	return ""
}

func trimDigest(ref string) string {
	if idx := strings.Index(ref, "@"); idx != -1 {
		return ref[:idx]
	}
	return ref
}
