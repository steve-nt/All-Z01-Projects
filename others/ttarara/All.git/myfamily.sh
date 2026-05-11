#!/bin/bash

# Retrieve the JSON data from the URL using curl and filter it with jq based on HERO_ID
family=$(curl -s https://platform.zone01.gr/assets/superhero/all.json | jq '.[] | select(.id == '$HERO_ID') | .connections.relatives')


family=$(echo "$family" | tr -d '"')

echo "$family"
