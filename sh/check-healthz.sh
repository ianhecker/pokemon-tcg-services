#!/bin/sh

SERVICES="localhost:8080/healthz"

for URL in $SERVICES; do
    echo "Checking $url..."
    if curl -fs "$URL" >/dev/null 2>&1; then
        echo "✅ $URL"
    else
        echo "❌ $URL"
        exit 1
    fi
done