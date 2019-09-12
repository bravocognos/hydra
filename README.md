# git-submodule-concurrent
Concurrently runs commands against git submodules

## Overview
Bored of running syncronous commands for git submodules? Meet `gsc`.

## Usage
`$ gsc "npm install"`

## Build
1. Clone the project
2. Run `go build -o gsc`

### Move to the bin folder
Move the binary to any folder that's included in your `$PATH`

#### Mac OS X
1. Run `mv gsc /usr/local/bin/`
