#!/usr/bin/env bash

[ -z "$MOCKERY_VERSION" ] && MOCKERY_VERSION="2.36.0"

install_mockery () {
  case "$2" in
      Darwin*)
        {
          PLATFORM=$2
          HARDWARE=$(uname -m)
        };;
      Linux*)
        {
          PLATFORM=$2
          HARDWARE=$(uname -m)
           if [ "$HARDWARE" == "aarch64" ]; then \
                HARDWARE="arm64"
              fi
        };;
  esac

  mkdir -p /tmp/mockery-"$1"

  curl -sL https://github.com/vektra/mockery/releases/download/v"$1"/mockery_"$1"_"$PLATFORM"_"$HARDWARE".tar.gz | tar xvz -C /tmp/mockery-"$1" \
    && mv /tmp/mockery-"$1"/mockery "$GOPATH"/bin/mockery
}

osType="$(uname -s)"

case "${osType}" in
    Darwin*)
      {
        # checking if mockery is installed
        if brew ls --versions mockery > /dev/null; then
          echo ">> uninstalling mockery via brew... "; \
          echo ">> brew uninstall mockery "; \
          exit 1
        fi
      };;
    Linux*)
      {
        # checking if mockery is installed
        if dpkg -l | grep mockery > /dev/null; then
          echo ">> uninstalling mockery via apt-get... "; \
          echo ">> apt-get remove -y mockery "; \
          exit 1
        fi
      };;
    *)
      {
        echo "Unsupported OS, exiting"
        exit 1
      } ;;
esac

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# adding GOBIN to PATH
[[ ":$PATH:" != *"$GOPATH/bin"* ]] && PATH=$PATH:"$GOPATH"/bin

# checking if mockery is available and it is the version specify
if ! command -v mockery > /dev/null; then \
    echo ">> Installing mockery v$MOCKERY_VERSION...";
    install_mockery "$MOCKERY_VERSION" "$osType"
else
  VERSION_INSTALLED="$(mockery --version --quiet | cut -d' ' -f2)"
  if [ "${VERSION_INSTALLED}" != "v${MOCKERY_VERSION}" ]; then \
    echo ">> Updating mockery form "${VERSION_INSTALLED}" to v$MOCKERY_VERSION..."; \
    install_mockery "$MOCKERY_VERSION" "$osType"
  fi
fi
