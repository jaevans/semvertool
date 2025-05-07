# semvertool

`semvertool` is a CLI tool to manage semantic versioning (semver) strings. It provides commands to bump major, minor, patch, and prerelease versions, as well as additional utilities for working with semver strings.

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

### bump git

Bump a version based on the latest semver tag in the git repository.

```shell
Examples:
git tag v0.1.0
semvertool bump git
v0.1.1

git tag v1.0.0
semvertool bump git --minor
v1.1.0

git tag v1.0.0-alpha.0
semvertool bump git --prerelease
v1.0.0-alpha.1

git tag v1.0.0-alpha.1
semvertool bump git --prerelease --hash
v1.0.0-alpha.2+3f6d1270

git tag v1.0.0-alpha.2+3f6d1270
semvertool bump git --minor
v1.1.0
```

### script

Provides utilities for scripting with semantic versions. These commands are designed to be used in shell scripts, returning exit codes that can be used in conditionals.

#### script compare

Compare two semantic versions and return an exit code based on the comparison.

```shell
semvertool script compare <version1> <version2>
```

Exit codes:

- **0**: When versions are equal
- **11**: When version1 is greater than version2 (version1 is newer)
- **12**: When version1 is less than version2 (version2 is newer)

Examples:

```shell
semvertool script compare 1.0.0 1.0.0
echo $? # Returns 0 (equal)

semvertool script compare 2.0.0 1.0.0
echo $? # Returns 11 (first version is newer)

semvertool script compare 1.0.0 1.1.0
echo $? # Returns 12 (second version is newer)
```

#### script released

Check if a version is a release version (not a prerelease and has no metadata).

```shell
semvertool script released <version>
```

Exit codes:

- **0**: If the version is a release version (X.Y.Z only)
- **1**: If the version is a prerelease or has metadata

Examples:

```shell
semvertool script released 1.0.0
echo $? # Returns 0 (it's a release version)

semvertool script released 1.0.0-alpha.1
echo $? # Returns 1 (it's a prerelease)

semvertool script released 1.0.0+build.123
echo $? # Returns 1 (it has metadata)
```

### `sort`

Sorts a list of semver strings in ascending order. This is useful for organizing version lists or ensuring proper version ordering.

### Usage

```bash
semvertool sort <version1> <version2> <version3> ...
```

Example:

```bash
semvertool sort 1.0.0 2.0.0 0.9.0
0.9.0 1.0.0 2.0.0

semvertool sort --order descending 1.0.0 2.0.0 0.9.0
2.0.0 1.0.0 0.9.0

semvertool sort  1.0.0 1.0.0-alpha.1 1.0.0-alpha.2
1.0.0-alpha.1 1.0.0-alpha.2 1.0.0

semvertool sort --no-prerelease 1.0.0 1.0.1-alpha.1 2.0.0 1.0.1 1.0.1-beta.1
1.0.0 1.0.1 2.0.0
```

### previous

Get the previous semver tag from git history. This is useful for determining what version preceded the current one.

```shell
semvertool previous
```

The command works in two modes:

1. If the current commit has a semver tag, it will return the tag that came before it
2. If the current commit is not tagged, it will return the previous tag in the history

The command also provides a `--released` flag that filters out prerelease versions:

```shell
semvertool previous --released
```

Examples:

For a repository with tags `v1.0.0`, `v1.1.0`, `v1.2.0-alpha.1`, `v1.2.0`:

```shell
# When HEAD is at v1.2.0
semvertool previous
v1.2.0-alpha.1

# Same scenario, but only looking at released versions
semvertool previous --released
v1.1.0

# When HEAD is at an untagged commit after v1.2.0
semvertool previous
v1.2.0

# When there's only one tag in the repository
semvertool previous
Error: no previous tag available - already at oldest tag

# Get the previous tag from a different repository
semvertool previous -r /path/to/other/git/repo

# Short form of the repository flag
semvertool previous --repository=/path/to/other/git/repo
```
