package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdaterUpdateMatchingImage(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "updater-docker-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "Dockerfile")
	original := "FROM ghcr.io/semrels/app:1.2.3\nRUN echo hi\n"
	if err := os.WriteFile(file, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := NewUpdater().Update(file, "semrels/app", "ghcr.io", "1.3.0"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "FROM ghcr.io/semrels/app:1.3.0") {
		t.Fatalf("updated file = %s", got)
	}
}

func TestUpdaterMissingFile(t *testing.T) {
	t.Parallel()

	err := NewUpdater().Update(filepath.Join(t.TempDir(), "Dockerfile"), "app", "", "1.3.0")
	if err == nil || !strings.Contains(err.Error(), "read") {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestUpdaterImageNotFound(t *testing.T) {
	t.Parallel()

	updated, changed, err := updateContent("FROM alpine:3.19\n", "busybox", "", "1.3.0")
	if err != nil {
		t.Fatalf("updateContent() error = %v", err)
	}
	if changed || updated == "" {
		t.Fatalf("expected no change, got changed=%v content=%q", changed, updated)
	}
}

func TestTrimRegistryAndTag(t *testing.T) {
	t.Parallel()

	if got := trimRegistryAndTag("ghcr.io/semrels/app:1.2.3"); got != "semrels/app" {
		t.Fatalf("trimRegistryAndTag() = %q", got)
	}
}

func TestUpdateContentUsesFirstImageWhenUnset(t *testing.T) {
	t.Parallel()

	updated, changed, err := updateContent("FROM alpine:3.19\n", "", "", "3.20")
	if err != nil {
		t.Fatalf("updateContent() error = %v", err)
	}
	if !changed || !strings.Contains(updated, "FROM alpine:3.20") {
		t.Fatalf("updated=%q changed=%v", updated, changed)
	}
}

func TestBuildReferenceWithRegistry(t *testing.T) {
	t.Parallel()

	got := buildReference("ghcr.io/semrels/app:1.2.3", "semrels/app", "docker.io", "1.3.0")
	if got != "docker.io/semrels/app:1.3.0" {
		t.Fatalf("buildReference() = %q", got)
	}
}

func TestTrimHelpers(t *testing.T) {
	t.Parallel()

	if got := trimTag("ghcr.io/semrels/app:1.2.3"); got != "1.2.3" {
		t.Fatalf("trimTag() = %q", got)
	}
	if got := trimDigest("alpine@sha256:abc"); got != "alpine" {
		t.Fatalf("trimDigest() = %q", got)
	}
}
