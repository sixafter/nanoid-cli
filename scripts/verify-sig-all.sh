#!/usr/bin/env bash
# Copyright (c) 2024-2025 Six After
set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}/os-type.sh"

if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

curl_retry() {
    local url="$1"
    local out="$2"
    local attempt=1
    local max=5
    local delay=2

    while true; do
        # -f: fail on HTTP error codes
        # -s: silent
        # -S: show errors
        # -L: follow redirects
        if curl -fSsSL "${url}" -o "${out}"; then
            return 0
        fi

        if (( attempt >= max )); then
            echo "[ERROR] curl failed after ${attempt} attempts: ${url}" >&2
            return 1
        fi

        echo "[WARN] curl failed (attempt ${attempt}/${max}). Retrying in ${delay}s..."
        sleep $delay
        attempt=$(( attempt + 1 ))
        delay=$(( delay * 2 ))  # exponential backoff
    done
}

PROJECT="nanoid-cli"
REPO="sixafter/${PROJECT}"

TMP="${__dir}/tmp"
mkdir -p "${TMP}"

echo "Project: $PROJECT"
echo "Repository: $REPO"
echo "Artifact directory: $TMP"
echo

TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | jq -r .tag_name)
VERSION=${TAG#v}

echo "Latest release: ${TAG} (version ${VERSION})"
echo

if command -v sha256sum >/dev/null 2>&1; then
    SHA256="sha256sum"
else
    SHA256="shasum -a 256"
fi

echo "Fetching asset list..."
ASSETS=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | jq -r '.assets[].name')

echo
echo "Downloading all release artifacts..."

for asset in $ASSETS; do
  echo "→ $asset"
  curl_retry \
    "https://github.com/${REPO}/releases/download/${TAG}/${asset}" \
    "${TMP}/${asset}"
done

echo
echo "Downloading checksums.txt and signature..."
curl_retry \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt" \
  "${TMP}/checksums.txt"

curl_retry \
  "https://github.com/${REPO}/releases/download/${TAG}/checksums.txt.sigstore.json" \
  "${TMP}/checksums.txt.sigstore.json"

echo
echo "Verifying checksums.txt signature..."
cosign verify-blob \
  --key "${__dir}/../cosign.pub" \
  --bundle "${TMP}/checksums.txt.sigstore.json" \
  "${TMP}/checksums.txt"

echo "✓ checksums.txt signature OK"

echo
echo "Verifying artifact signatures..."

for asset in $ASSETS; do
    # Skip SBOM files—they have no associated signature
    if [[ "$asset" == *.sbom.json ]]; then
        continue
    fi

    bundle="${asset}.sigstore.json"
    if [[ ! -f "${TMP}/${bundle}" ]]; then
        echo "[ERROR] Missing signature bundle for: ${asset}"
        exit 1
    fi

    echo "→ Verifying signature for $asset"
    cosign verify-blob \
      --key "${__dir}/../cosign.pub" \
      --bundle "${TMP}/${bundle}" \
      "${TMP}/${asset}"

    echo "   ✓ OK"
done

echo
echo "Verifying checksums for all artifacts..."
(
    cd "${TMP}"
    $SHA256 -c checksums.txt
)

echo
echo "✔ ALL signatures and checksums verified successfully."
