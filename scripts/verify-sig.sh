#!/usr/bin/env bash
# Copyright (c) 2024-2025 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.

set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}/os-type.sh"

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

# ------------------------------------------------------------
# Project / repository name (portable)
# ------------------------------------------------------------
PROJECT="nanoid-cli"
REPO="sixafter/${PROJECT}"
MODULE="github.com/${REPO}"

# tmp directory for artifacts
TMP="${__dir}/tmp"
mkdir -p "${TMP}"

echo "Project: ${PROJECT}"
echo "Repository: ${REPO}"
echo "Module path: ${MODULE}"
echo "Artifact directory: ${TMP}"
echo

# ------------------------------------------------------------
# Detect latest release
# ------------------------------------------------------------
TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | jq -r .tag_name)
VERSION=${TAG#v}

echo "Latest release: ${TAG} (version: ${VERSION})"

# ------------------------------------------------------------
# Compute Go-style OS/ARCH for artifact naming
# ------------------------------------------------------------
GOOS=$(goos)
GOARCH=$(goarch)

if [[ "$GOOS" == "unsupported" ]] || [[ "$GOARCH" == "unsupported" ]]; then
    echo "[ERROR] Unsupported OS/ARCH: ${GOOS}/${GOARCH}"
    exit 1
fi

ARTIFACT="nanoid_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
echo "Using artifact: ${ARTIFACT}"
echo

# ------------------------------------------------------------
# Determine SHA-256 tool
# ------------------------------------------------------------
if command -v sha256sum >/dev/null 2>&1; then
  SHA256="sha256sum"
else
  SHA256="shasum -a 256"
fi

# ------------------------------------------------------------
# Download release artifacts → tmp/
# ------------------------------------------------------------
echo
echo "Downloading release artifacts into ${TMP}..."

curl -sSfL -o "${TMP}/${ARTIFACT}" \
  "https://github.com/${REPO}/releases/download/${TAG}/${ARTIFACT}"

curl -sSfL -o "${TMP}/${ARTIFACT}.sigstore.json" \
  "https://github.com/${REPO}/releases/download/${TAG}/${ARTIFACT}.sigstore.json"

curl -sSfL -o "${TMP}/checksums.txt" \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt"

curl -sSfL -o "${TMP}/checksums.txt.sigstore.json" \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt.sigstore.json"

# ------------------------------------------------------------
# Verify artifact signature
# ------------------------------------------------------------
echo
echo "Verifying artifact signature..."

cosign verify-blob \
  --key "${__dir}/../cosign.pub" \
  --bundle "${TMP}/${ARTIFACT}.sigstore.json" \
  "${TMP}/${ARTIFACT}"

echo "Artifact signature OK."

# ------------------------------------------------------------
# Verify checksums manifest signature
# ------------------------------------------------------------
echo
echo "Verifying checksums.txt signature..."

cosign verify-blob \
  --key "${__dir}/../cosign.pub" \
  --bundle "${TMP}/checksums.txt.sigstore.json" \
  "${TMP}/checksums.txt"

echo "Checksums signature OK."

# ------------------------------------------------------------
# Validate local artifact integrity
# ------------------------------------------------------------
echo
echo "Verifying file checksum for ${ARTIFACT}..."

(
  cd "${TMP}"

  # extract only the line relating to our artifact
  LINE=$(grep "${ARTIFACT}" checksums.txt)

  if [[ -z "${LINE}" ]]; then
    echo "❌ No checksum entry found for ${ARTIFACT} in checksums.txt"
    exit 1
  fi

  echo "${LINE}" | $SHA256 -c -
) || {
    echo
    echo "❌ Release verification FAILED."
    exit 1
}

echo
echo "✔ Release verification succeeded."
