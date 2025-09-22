#!/bin/sh

TCG_PLAYER_PRODUCT_ID=106999

curl "https://api.justtcg.com/v1/cards?tcgplayerId=$TCG_PLAYER_PRODUCT_ID" \
    -H "x-api-key:$JUST_TCG_API_KEY" | jq .