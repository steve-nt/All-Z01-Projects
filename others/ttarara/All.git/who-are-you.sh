url="https://platform.zone01.gr/assets/superhero/all.json"

name=$(curl -s "$url" | jq -r '.[] | select(.id == 70) | .name')

# Print the name surrounded by double quotes and followed by a newline
printf "\"%s\"\n" "$name"
