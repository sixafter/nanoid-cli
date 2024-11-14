# NanoID CLI

A simple, fast, and concurrent command-line tool for generating secure, URL-friendly unique string IDs 
using the [NanoID](https://github.com/sixafter/nanoid) Go implementation.

[![CI](https://github.com/sixafter/nanoid-cli/workflows/ci/badge.svg)](https://github.com/sixafter/nanoid-cli/actions)
[![Go](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid-cli)](https://img.shields.io/github/go-mod/go-version/sixafter/nanoid-cli)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=six-after_nano-id-cli&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=six-after_nano-id-cli)
[![GitHub issues](https://img.shields.io/github/issues/sixafter/nanoid-cli)](https://github.com/sixafter/nanoid-cli/issues)
[![Go Reference](https://pkg.go.dev/badge/github.com/sixafter/nanoid-cli.svg)](https://pkg.go.dev/github.com/sixafter/nanoid-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/sixafter/nanoid-cli)](https://goreportcard.com/report/github.com/sixafter/nanoid-cli)
[![License: Apache 2.0](https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square)](LICENSE)

## Features

- **Customizable Length**: Specify the length of the generated Nano ID.
- **Custom Alphabet**: Define your own set of characters for ID generation.
- **Multiple ID Generation**: Generate multiple IDs in a single command.
- **Output Options**: Print IDs to the console or save them to a file.
- **Verbose Mode**: Enable detailed logs during ID generation.

## Installation

### Clone the repository and build the executable:

```sh
git clone https://github.com/sixafter/nanoid-cli.git
cd nanoid-cli
make build
```

This command compiles the main.go file and produces an executable named `nanoid` in the `./out` directory.

### Homebrew

```sh
brew tap sixafter/nanoid
brew install nanoid
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
./nanoid generate --alphabet "abcdef123456"
```

Output:

```sh
1a2b3c4d5e6f1a2b3c4d5e6f1a2b3c4
```

Generate Multiple Nano IDs and Save to a File:

```sh
nanoid generate --count 5 --output ids.txt
```

Contents of `ids.txt`:

```sh
V1StGXR8_Z5jdHi6B-myT
mJzY8fK3Lq7B9sR2dT4hV
aB3cD4E5f6G7hI8jK9lMn
O1pQ2rS3tU4vW5xY6zA7b
D8eF9gH0iJ1kL2mN3oP4q
```

---

## Contributing

Contributions are welcome. See [CONTRIBUTING](CONTRIBUTING.md)

---

## License

This project is licensed under the [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/). See [LICENSE](LICENSE) file.
