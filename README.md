# finas

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/toozej/finas)
[![Go Report Card](https://goreportcard.com/badge/github.com/toozej/finas)](https://goreportcard.com/report/github.com/toozej/finas)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/toozej/finas/cicd.yaml)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/toozej/finas/total)

FINAS Is Not A Shell

## About
FINAS is an opinionated, expandable CLI written in Go which allows for execution of canned `docker run` commands. FINAS uses Viper + Cobra to generate sub-commands from loaded JSON config files stored in standard directories, and is heavily based on https://github.com/toozej/golang-starter as its project scaffolding.

## Usage

1. Copy or create some canned Docker commands
```bash
mkdir ~/.config/finas && cp configs/*.json ~/.config/finas/ 
```

2. List available commands with `finas --help`

## Development

Development tasks such as building, testing, and releasing are all `make` driven. Type `make` to see a list of available tasks.

### changes required to update golang version
- `make update-golang-version` 

## Inspiration
[GNU](https://www.gnu.org/gnu/about-gnu.html) + [thefuck](https://github.com/nvbn/thefuck) = FINAS, or `f` for short.
This name was chosen largely because I can never remember the arguments for some handy Docker run commands, and I wanted a quick and easy way to run them to avoid yelling expletives at my computer.
