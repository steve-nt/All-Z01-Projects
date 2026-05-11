url="https://platform.zone01.gr/assets/superhero/all.json"
name=$(curl -s "$url" | jq -r '.[] | select(.id == 170) | .name')
power=$(curl -s "$url" | jq -r '.[] | select(.id == 170) | .powerstats | .power')
gender=$(curl -s "$url" | jq -r '.[] | select(.id == 170) | .appearance | .gender')
# Printing the extracted values

printf "%s\n" "$name"
printf "%s\n" "$power"
printf "%s\n" "$gender"