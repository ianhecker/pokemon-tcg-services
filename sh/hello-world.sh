#!/bin/sh

# Unique ID for 1st Edition Shadowless Base Set Charizard
# https://www.tcgplayer.com/product/106999/pokemon-base-set-shadowless-charizard
TCG_PLAYER_PRODUCT_ID=106999

curl "localhost:8080/v1/cards?id=$TCG_PLAYER_PRODUCT_ID" | jq .
