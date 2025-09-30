# Pokemon TCG Services

Microservices for collecting trading cards from the Pokemon Trading Card Game (TCG)

<div align="center">
  <img src="assets/logo.png" alt="My mascot" height="400"/>
</div>

[![Test](https://github.com/ianhecker/pokemon-tcg-services/actions/workflows/test.yml/badge.svg)](https://github.com/ianhecker/pokemon-tcg-services/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ianhecker/pokemon-tcg-services)](https://goreportcard.com/report/github.com/ianhecker/pokemon-tcg-services)
![Test Coverage](https://img.shields.io/endpoint?url=https%3A%2F%2Fianhecker.github.io%2Fpokemon-tcg-services%2Fcoverage.json)

## Table of Contents
- [Overview](#overview)
    - [Pricing Standards](#pricing-standards)
    - [The Mint Condition Problem](#the-mint-condition-problem)
    - [Used Card Comps](#used-card-comps)
    - [Comps Across Conditions](#comps-across-conditions)
    - [Project Goal](#project-goal)
- [Architecture](#architecture)
    - [Card Pricer](#card-pricer)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
    - [Testing](#testing)
    - [Running](#running)

---

## Overview

The Pokémon Trading Card Game `TCG` released its first set in the US in 1998

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
- **Mint** - `M`
- **Near Mint** - `NM`
- **Lightly Played** - `LP`
- **Moderately Played** - `MP`
- **Heavily Played** - `HP`
- **Damaged** - `D`

Even with standardized conditions, used cards are often ungraded, and their
condition is negotiable

As a result, estimating value requires checking `comps` across multiple
conditions

### Comps Across Conditions

Some tools can `comp` for lower conditions, but only one at a time

Checking multiple conditions means repeated searches and visual parsing

No aggregated price history is stored

### Project Goal

This project simplifies checking `comps`

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

> Docker images are orchestrated by Kubernetes

Coming soon!

Services
- `card-pricer`

### Card Pricer

For a card, for all printings, for all conditions, fetches *all available* prices

Using a unique TCG Player Product ID, queries Just TCG for card prices for all
conditions

---

## Installation

See [INSTALL.md](./docs/INSTALL.md)

---

## Configuration

See [CONFIGURATION.md](./docs/CONFIGURATION.md)

---

## Usage

[Makefile](./Makefile) is your one-stop-shop for all repository commands

### Testing

See [TEST.md](./docs/TEST.md) for testing

See [MOCKS.md](./docs/MOCKS.md) for mocks and mock generation

### Running

See [RUN.md](./docs/RUN.md)

---
