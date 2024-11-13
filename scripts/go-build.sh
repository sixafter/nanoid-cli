#!/bin/bash
# Copyright (c) 2024 Six After, Inc
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

set -e

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

source "${__dir}"/os-type.sh

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

X_BUILD_VERSION=${GitVersion_AssemblySemFileVer:-0.1.0}
echo "[INFO] Build version set to '${X_BUILD_VERSION}'." >&1

X_BUILD_COMMIT=$(git rev-parse HEAD 2> /dev/null || true)
echo "[INFO] Build commit set to '${X_BUILD_COMMIT}'." >&1

export LDFLAGS="-s -w -X github.com/sixafter/nanoid-cli/cmd/version.version=${X_BUILD_VERSION} -X github.com/sixafter/nanoid-cli/cmd/version.gitCommitID=${X_BUILD_COMMIT}"
go build -o "${BINARY_NAME}" -ldflags "${LDFLAGS}" main.go
