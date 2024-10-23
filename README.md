# semvertool

## Description

A small toolkit for managing semantic versioning in your projects.

## Features

- Full support for semantic versioning
- Bump major, minor, patch, or pre-release versions
- Embed build metadata in the version string
- Select the type of bump from a text string (git commit message)
- Compare two versions (TBD)
- Validate a version string (TBD)

## Installation

## From Source

To install this application from the git repository, you can follow these steps:

1. Clone the repository: `git clone https://github.com/jaevans/semvertool.git`
2. Navigate to the project directory: `cd semvertool`
3. Install the dependencies: `go mod download`
4. Build the application: `go build -o semvertool`

## Usage

Here is a list of the available commands:

### bump

```shell
semvertool bump --major 0.1.0-alpha.1+build.1
1.0.0

semvertool bump --minor 1.0.0
1.1.0

semvertool bump --patch 1.1.0
1.1.1

semvertool bump --prerelease 1.1.1-alpha.0
1.1.1-alpha.1

semvertool bump --prerelease 1.1.1-alpha.1.0
1.1.1-alpha.1.1

semvertool bump --build $(git rev-parse --short HEAD) 1.1.1-alpha.1
1.1.1-alpha.2+3f6d1270

semvertool bump 1.1.1-alpha.2+3f6d1270
```

### git

Bump a version based on the latest semver tag in the git repository.

```shell
Examples:
git tag v0.1.0
semvertool git
v0.1.1

git tag v1.0.0
semvertool git --minor
v1.1.0

git tag v1.0.0-alpha.0
semvertool git --prerelease
v1.0.0-alpha.1

git tag v1.0.0-alpha.1
semvertool git --prerelease --hash
v1.0.0-alpha.2+3f6d1270

git tag v1.0.0-alpha.2+3f6d1270
semvertool git --minor
v1.1.0
```

### compare

Compare two semver versions and return an exit code based on their comparison.

```shell
Examples:
semvertool compare 1.0.0 2.0.0
2

semvertool compare 2.0.0 2.0.0
1

semvertool compare 2.0.0 1.0.0
0
```
