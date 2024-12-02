#!/usr/bin/env bash

[ -z "$DOCKER_IMAGE_TAG" ] && DOCKER_IMAGE_TAG="latest"
[ -z "$DOCKERFILE_PATH" ] && DOCKERFILE_PATH="Dockerfile"
[ -z "$DOCKERBUILD_CONTEXT" ] && DOCKERBUILD_CONTEXT="."
[ -z "$DOCKER_SECRET" ] && DOCKER_SECRET=false

# Getting name of the project
# gh_repo is the variable available in gh-repo.sh that contains the long repo name. Example: github.com/consensys-vertical-apps/platform-data-pipeline-toolkit
BASEDIR=$(dirname "$0")
source "$BASEDIR"/gh-repo.sh

gh_repo=${gh_repo#"github.com/"}
gh_repo=$(echo "$gh_repo" | awk '{print tolower($0)}')

echo "Building docker image ..."

if [ ! -z "$DOCKER_GITHUB_TOKEN" ]; then \
  if [[ "$DOCKER_SECRET" == true ]]; then \
    # if $DOCKER_GITHUB_TOKEN AND $DOCKER_SECRET exists then build with secret
    GH_ACCESS_TOKEN="$DOCKER_GITHUB_TOKEN" docker buildx build -t "$gh_repo:$DOCKER_IMAGE_TAG" -f "$DOCKERFILE_PATH" --secret id=GH_ACCESS_TOKEN "$DOCKERBUILD_CONTEXT"
  else
    docker buildx build -t "$gh_repo:$DOCKER_IMAGE_TAG" -f "$DOCKERFILE_PATH" --build-arg GH_ACCESS_TOKEN="$DOCKER_GITHUB_TOKEN" "$DOCKERBUILD_CONTEXT"
  fi
else
  docker build -t "$gh_repo:$DOCKER_IMAGE_TAG" -f "$DOCKERFILE_PATH" "$DOCKERBUILD_CONTEXT"
fi