# updater-docker

Docker image updater plugin for SemRel.

Builds Docker images and publishes tagged releases to Docker Hub as part of SemRel automation.

## Documentation

- SemRel docs (planned): <https://github.com/SemRels/semrel/tree/main/docs/plugins/updater-docker>
- Plugin template: <https://github.com/SemRels/plugin-template>
- Registry: <https://registry.semrel.io>

## Repository Layout

~~~text
cmd/plugin/              Plugin entry point
internal/plugin/         Business logic scaffold
internal/grpc/           gRPC transport scaffold
proto/v1                 Symlink to the SemRel protobuf contract
.github/workflows/       CI, release, and security automation
~~~

## Development

~~~bash
go build ./cmd/plugin
go test ./...
~~~

## Configuration Example

~~~yaml
plugins:
  - name: updater-docker
    type: updater
    config:
      dockerfile: Dockerfile
      image: semrels/example
      tags:
        - latest
        - ${version}
      registry: docker.io
      username: ${DOCKERHUB_USERNAME}
      password: ${DOCKERHUB_TOKEN}
~~~

## Status

This repository is bootstrapped from SemRels/plugin-template and is ready for implementation.
