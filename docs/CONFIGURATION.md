# Configuration

All configuration is done with environment variables

You will need the following:
1) [Just TCG API Key](https://justtcg.com/auth/signup)

Docker images require a configuration file. See [Configuration Files](#configuration-files)

## Environment Variables

```bash
export JUST_TCG_API_KEY=
```

## Configuration Files

To get environment variables into Docker images easily, a `.env` configuration
file is passed when running

```bash
docker run ... --env-file .env ...
```

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
