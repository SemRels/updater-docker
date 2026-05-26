# updater-docker

Updates a version argument in a Dockerfile.

This plugin is distributed as the standalone Go binary `semrel-plugin-updater-docker`. Semrel executes the binary as a subprocess, provides plugin configuration through `SEMREL_PLUGIN_*` environment variables, provides release context through `SEMREL_*` environment variables, reads standard output, and treats exit code `0` as success and any non-zero exit code as failure. Install the binary in `~/.semrel/plugins/` or anywhere on your `$PATH`.

## Installation

```bash
go install github.com/SemRels/updater-docker/cmd/plugin@latest
```

## Configuration

```yaml
plugins:
  - name: updater-docker
    path: ~/.semrel/plugins/semrel-plugin-updater-docker
    env:
      SEMREL_PLUGIN_FILE: "Dockerfile"
      SEMREL_PLUGIN_ARG_NAME: "VERSION"
```

## `SEMREL_PLUGIN_*` variables

| Name | Required | Description | Default |
| --- | --- | --- | --- |
| `SEMREL_PLUGIN_FILE` | Optional | Path to the Dockerfile to update. | Dockerfile |
| `SEMREL_PLUGIN_ARG_NAME` | Optional | Docker build argument name that stores the version. | VERSION |

## `SEMREL_*` release context used

| Variable | Description |
| --- | --- |
| `SEMREL_VERSION` | Resolved release version for the current run. |
| `SEMREL_NEXT_VERSION` | Next version computed by semrel for the release. |
| `SEMREL_DRY_RUN` | Whether semrel is running in dry-run mode. |

## Example behavior

The plugin updates the configured Docker build argument to the new release version and prints what changed.

## License

Apache-2.0
