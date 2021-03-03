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

## License

This project is [licensed](./LICENSE) under Apache 2.0.
