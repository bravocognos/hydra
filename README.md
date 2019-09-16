[![golang version status](https://img.shields.io/badge/Golang-1.12.x-green.svg?style=flat-square)]()

![Banner](/resources/banner.png?raw=true "Banner")

# Hydra
Concurrently runs commands against git submodules.

## Overview
Bored of running syncronous commands for git submodules? Meet `hydra`.

## Usage
Download from the release page, move to a directory included in your `$PATH`, and and run. Few examples:

- `$ hydra "npm install"`
- `$ hydra "rm -rdf node_modules"`
- `$ hydra "git pull"`

## Developing
You will need Golang. It was developed using version `1.12.9`

1. Clone the project
2. Make changes
3. Build `$ go build -o hydra`
4. Test
5. Open a Merge Request referecing the opened issues: `Ref.: #ISSUE_NUMBER`
