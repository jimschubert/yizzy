# yizzy

YAML migrations using yq and selector syntax/operators. It can simplify DevOps administration of CI files written in YAML.

* _Need to upgrade hundreds of CircleCI configs with complex logic?_ - **Get yizzy**.
* _Want apply multiple changes consistently to all GitHub Workflows on your machine?_ - **Get yizzy**.
* _Prefer awk/sed and bash scripting?_ - **Don't get yizzy**.

[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/yizzy)
![Go](https://github.com/jimschubert/yizzy/workflows/Build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/yizzy)](https://goreportcard.com/report/github.com/jimschubert/yizzy)
![Docker Pulls](https://img.shields.io/docker/pulls/jimschubert/yizzy)
<!-- [![codecov](https://codecov.io/gh/jimschubert/yizzy/branch/master/graph/badge.svg)](https://codecov.io/gh/jimschubert/yizzy) --> 

## Installation

### Binaries

Latest binary releases are available via [GitHub Releases](https://github.com/jimschubert/yizzy/releases).

### Homebrew

If you prefer homebrew, I got you.

```
brew install jimschubert/tap/yizzy
```

or:

```
brew tap jimschubert/tap
brew install yizzy
```

## Usage

```shell
$ yizzy -h
Usage:
  yizzy [OPTIONS]

Application Options:
  -f, --file=FILE    The file to process
  -d, --dir=         The directory where migrations reside
      --in-place     Writes a file in place
  -v, --version      Display version information

Help Options:
  -h, --help         Show this help message
```

### Migration document(s)

A migration document holds the list of operations and environments we intend to apply to a target file. The document applies one or more 
[expressions as supported by yq](https://mikefarah.gitbook.io/yq/) against a YAML document. 

####  Migration: `operations`

An operation consists of:

* `value_type`: an optional, defining the [YAML type](https://yaml.org/type/) which will be applied by `value` or `eval`
* `value`: a scalar or literal value which does not require document traversal or evaluation
* `eval`: an expression requiring traversal or [yq operations](https://mikefarah.gitbook.io/yq/operators) evaluated in the
  context defined by `selector` (or the document root by default)
* `selector`: a [yq](https://mikefarah.gitbook.io/yq/) expression targeting one or more nodes on which to operate contextually
  when applying the `value` or `eval` result
  
#### Migration: `env`

The `env` node is a simple map of keys representing the environment variable to be set during evaluation of each operation,
and an expression to apply whenever that environment variable is referenced via yq's [env variable operators](https://mikefarah.gitbook.io/yq/operators/env-variable-operators).

**NOTE** A literal reference defined within the YAML document must result in an expression wrapped in quotes. In YAML, this
will be the following syntax:

```
LITERAL_ENV: '"Jim Schubert"'
```

This creates an expression of `"Jim Schubert"`. One set of single/double quotes would result in an unquoted expression `Jim Schubert`,
which would cause a parser error in yq.

#### Example

Create a directory called `migrations`, and create a dated YAML file within this directory. For example: `2021-03-05.yml`:

```yaml
env:
  FIRST_NAME: .bill-to.given
  LAST_NAME: .bill-to.family
operations:
  - selector: .ship-to
    eval: '.full_name = strenv(FIRST_NAME) + " " + strenv(LAST_NAME)'
    value_type: '!!str'
```

A migration document contains an optional map named `env`, a key/value mapping of environment variable names which will
be passed to each operation, and a selector evaluated against the input document prior to operations. These operations are
applied in declaration order.

See [testdata](./testdata) for examples of migrations and their expectations once applied to the YAML 1.2 specification's `invoice.yaml`.

## Contributors

### Build

Build a local distribution for evaluation using latest [goreleaser](https://goreleaser.com/).

```bash
goreleaser release --skip-publish --snapshot --rm-dist
```

This will create an executable application for your os/architecture under `dist`:

```
dist
├── yizzy_darwin_amd64
│   └── yizzy
├── yizzy_linux_386
│   └── yizzy
├── yizzy_linux_amd64
│   └── yizzy
├── yizzy_linux_arm64
│   └── yizzy
├── yizzy_linux_arm_6
│   └── yizzy
└── yizzy_windows_amd64
    └── yizzy.exe
```

Build and execute locally:

* Get dependencies
  ```shell
  go get -d ./...
  ```
* Build
  ```shell
  go build -o yizzy cmd/main.go
  ```
* Run
  ```shell
  ./yizzy --help
  ```
### Linting

1. Install [golangci-lint](https://golangci-lint.run/usage/install/#local-installation).
2. Run it: `golangci-lint run ./...`

## License

This project is [licensed](./LICENSE) under Apache 2.0.
