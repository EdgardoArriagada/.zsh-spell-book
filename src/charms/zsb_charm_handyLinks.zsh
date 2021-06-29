#!/usr/bin/env zsh
set -o errexit
set -o nounset
set -o pipefail

declare -r selection=$(print -l "${(@k)ZSB_CHARM_LINKS_DICTIONARY}" | rofi -show -dmenu -no-custom -i -p "Open link")

[[ -z "$selection" ]] && exit 0

xdg-open "${ZSB_CHARM_LINKS_DICTIONARY[$selection]}"

exit 0

# You must declare `ZSB_CHARM_LINKS_DICTIONARY` variable in ~/.zshenv
#
# example
#
# declare -Ar ZSB_CHARM_LINKS_DICTIONARY=(
#   [google]='https://www.google.com'
#   [runescape]='https://www.runescape.com'
# )
#
