#!/bin/bash
set -e

echo "Starting release process for version $CI_COMMIT_TAG..."

#? Verify this is a new version
if curl --silent "https://gitlab.com/api/v4/projects/$CI_PROJECT_ID/releases" | grep -q "\"tag_name\":\"$CI_COMMIT_TAG\""; then
  echo "❌ Release $CI_COMMIT_TAG already exists!"
  exit 1
fi

docker login "${CI_REGISTRY}" --username "${CI_REGISTRY_USER}" --password "${CI_REGISTRY_PASSWORD}"

if ! docker run --rm --privileged \
  -v "$PWD:/go/src/gitlab.com/goreleaser/cwc" \
  -w "/go/src/gitlab.com/goreleaser/cwc" \
  -v "/var/run/docker.sock:/var/run/docker.sock" \
  -e DOCKER_USERNAME="${CI_REGISTRY_USER}" \
  -e DOCKER_PASSWORD="${CI_REGISTRY_PASSWORD}" \
  -e DOCKER_REGISTRY="${CI_REGISTRY}" \
  -e GITLAB_TOKEN \
  goreleaser/goreleaser release --clean; then
  echo "❌ Release failed"
  exit 1
fi

echo "✅ Release completed successfully"
