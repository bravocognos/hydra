[![golang version status](https://img.shields.io/badge/Golang-1.12.x-green.svg?style=flat-square)]()

![Banner](/resources/banner.png?raw=true "Banner")

# Hydra
Concurrently runs commands against git submodules.

## Overview
Bored of running syncronous commands for git submodules? Meet `hydra`.

## Usage

1. Download from the release page
2. Move to a directory included in your `$PATH`, for example `/usr/local/bin` for Mac OS users.
3. Go to a repository with submodules
4. Run `hydra`. Few examples:

- `$ hydra "npm install"`
- `$ hydra "rm -rdf node_modules"`
- `$ hydra "git pull"`

## Developing
You will need Golang. It was developed using version `1.12.9`

1. Clone the project
2. Run `$ go get`
3. Make changes
4. Build `$ go build -o hydra`
5. Test your change(s)
6. Open a Merge Request referecing the opened issues: `Ref.: #ISSUE_NUMBER`
