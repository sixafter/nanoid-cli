# NanoID CLI

A simple, fast, and concurrent command-line tool for generating secure, URL-friendly unique string IDs 
using the [NanoID](https://github.com/sixafter/nanoid) Go implementation.

[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid-cli)](https://goreportcard.com/report/github.com/sixafter/nanoid-cli)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid-cli)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid-cli.svg)](https://pkg.go.dev/github.com/sixafter/nanoid-cli)

## Status

### Build & Test

[![CI](https://github.com/sixafter/nanoid-cli/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid-cli/actions)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid-cli)](https://github.com/sixafter/nanoid-cli/issues)

### Quality

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id-cli&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id-cli)
![CodeQL](https://github.com/sixafter/nanoid-cli/actions/workflows/codeql-analysis.yaml/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id-cli&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id-cli)

### Package and Deploy

[![Release](https://github.com/sixafter/nanoid-cli/workflows/release/badge.svg)](https://github.com/sixafter/nanoid-cli/actions)

## Features

- **Customizable Length**: Specify the length of the generated Nano ID.
- **Custom Alphabet**: Define your own set of characters for ID generation.
- **Multiple ID Generation**: Generate multiple IDs in a single command.
- **Verbose Mode**: Enable detailed logs during ID generation.

## Verify with Cosign

Download the binary and its `.sig` file, then run:

```sh
cosign verify-blob \
  --key https://raw.githubusercontent.com/sixafter/nanoid-cli/main/cosign.pub \
  --signature nanoid-cli-linux-amd64.tar.gz.sig \
  nanoid-cli-linux-amd64.tar.gz
```

If valid, Cosign will output:

```shell
Verified OK
```

## Installation

### Clone the repository and build the executable:

```sh
git clone https://github.com/sixafter/nanoid-cli.git
cd nanoid-cli
make build
```

This command compiles the `main.go` file and produces an executable named `nanoid` in the `./out` directory.

### Homebrew

```sh
brew tap sixafter/tap
brew install --cask nanoid
```

## Usage

Run the CLI to Generate a Default Nano ID:

```sh
nanoid generate
```

Output:

```sh
V1StGXR8_Z5jdHi6B-myT
```

Generate a Nano ID with a Custom Length:

```sh
nanoid generate --length 30
```

Output:

```sh
mJzY8fK3Lq7B9sR2dT4hV5nG1aC0eX
```

Generate a Nano ID with a Custom Alphabet:

```sh
nanoid generate --alphabet "abcdef123456"
```

Output:

```sh
1a2b3c4d5e6f1a2b3c4d5e6f1a2b3c4
```

Generate Multiple Nano IDs with verbose output:

```sh
nanoid generate --count 10 --verbose
```

Output:

```sh
_OKhyfsfINNfokJZxyj4j
HZUZ7sTHlLpub0rryyLsr
_agY2S55BoYSdipGVaL4P
FgrdoVAzzFZWS2bc42bre
saM2-PnvwIIyt312rkGbS
RCdECZCOr7VTkGXx5CoQo
bCX2GTzXJ22Azn0MAYkQ3
Fh7-65FYU9Higp7scLBht
uC87QtpSLb8ZX5oENCHJP
bTPg9AynQtzldZazM-wKV

Start Time..............: 2025-04-14T16:30:03-05:00
Total IDs generated.....: 10
Total time taken........: 46.959µs
Average time per ID.....: 4.695µs
Throughput..............: 212951.72 IDs/sec
Estimated output size...: 220 B
Estimated entropy per ID: 126.00 bits
Memory used.............: 0.32 MiB
```

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
