# Pokemon TCG Services

This repository aims to assist collectors obtain trading cards from the Pokemon
Trading Card Game (TCG)

---

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Services](#services)
- [Development](#development)

---

## Overview

The Pokémon Trading Card Game (TCG) released its first set in the US in 1998

Card releases span decades, and collectors must keep up with constantly changing
values

### Pricing Standards

[TCG Player](https://www.tcgplayer.com/) marketplace is the gold standard
for card pricing

Both collectors and businesses use it to look up `comps`, slang for comparable
listings, to estimate card values

### The Mint Condition Problem

Mainstream collecting heavily focuses on `Mint` and `Near Mint` conditions

Existing pricing tools are defaulted for those conditions

### Used Card Comps

The largest segment of the trading card market are cards below `Mint` conditions

Conditions:
- Mint - **M**
- Near Mint - **NM**
- Lightly Played - **LP**
- Moderately Played - **MP**
- Heavily Played - **HP**
- Damaged - **D**

Even with standardized conditions, used cards are often ungraded, and their
condition is negotiable

As a result, estimating value requires checking `comps` across multiple
conditions

### Comps Across Conditions

Some tools can comp for lower conditions, but only one at a time

Checking multiple conditions means repeated searches and visual parsing

No aggregated price history is stored

### Project Goal

This project simplifies checking comps

It removes the need for repeated searches with endless scrolling in the
TCG Player app.

For each card, it shows recent sales across all conditions in one view

Prices are current and aggregated from multiple marketplaces via
[Just TCG](https://justtcg.com/)

---

## Architecture

> All microservices are contained in one binary

This allows for easier packaging, distribution, and microservice execution

> Microservices are run through commands and flags

Docker images mimic those commands and flags for intuitive deployments

Services
- card-pricer

### Card Pricer

For a card, for all printings, for all conditions, fetches *all available* prices

Using a unique TCG Player Product ID, queries Just TCG for card prices for all
conditions

---

## Installation

### Required

- [Golang](https://go.dev/doc/install) - v1.24.2
- [Docker](https://docs.docker.com/desktop/) - To build images
    - [Ubuntu - Install via apt](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)

### Strongly Suggested

- [JQ](https://jqlang.org/) - For pretty JSON
- Makefile - For easy commands
- A Basic System Shell - For basic scripts in [sh](./sh)
    - Located at `/bin/sh`

### Optional

- [gotest](https://github.com/rakyll/gotest) - For colorized test output
    - `go install github.com/rakyll/gotest@latest`

---

## Configuration

### Environment Variables

> All configuration happens via environment variables

You will need to export:
- JUST_TCG_API_KEY = [Get a JUST TCG API Key](https://justtcg.com/dashboard/plans)

```bash
export JUST_TCG_API_KEY=
```

### Environment Configuration File

> You only need a `.env` file to run Docker containers

Services run in Docker containers get their environment variables when Docker
run is provided with a `.env` file, which it then injects as environment
variables into the container

Quickly create an `.env` from your already exported env variable
```bash
cat <<EOF > .env
JUST_TCG_API_KEY=${JUST_TCG_API_KEY}
EOF
```

Or, copy an empty template with `.env.example`
```bash
cp .env.example .env
```

---

## Usage

[Makefile](./Makefile) is your one-stop-shop for all repository commands

### Testing

```bash
make test
```

### Running

Build your binary
```bash
make build
```

Run the binary with the commands and flags you want
```bash
./bin/pokemon-tcg-services
```

Running card-pricer
```bash
 ./bin/pokemon-tcg-services service card-pricer --port=8080
```

You can skip all of these steps by simply running
```bash
make run-card-pricer
```

### Docker Images

```bash
make docker-build
```

### Hello World

Want to check if you can run services and fetch card prices?

*Terminal #1*
```bash
make run-card-pricer
```

*Terminal #2*
```bash
make hello-world
```

### Handy Scripts

See the [sh](./sh) directory for handy tools

Want to quickly fetch from the Just TCG API? See the [sh/justtcg/](./sh/justtcg/)
directory

---