# Run

Microservices each have respective commands and flags

See [CONFIGURATION.md](./CONFIGURATION.md) for required configuration before running

## Hello World

Want to check if you can run services?

*Terminal #1*
```bash
make run-card-pricer
```

*Terminal #2*
```bash
make hello-world
```

You should get pricing output for the Holo Rare Charizard Pokemon card from the
Shadowless Base Set!

## Card Pricer

Makefile
```bash
make run-card-pricer
```

Go run
```bash
go run main.go service card-pricer --port=8080
```

## Scripts

See the [sh](../sh) directory for handy scripting tools

Want to quickly fetch from the Just TCG API? See the [sh/justtcg/](../sh/justtcg/)
directory

## Kubernetes

Coming soon!