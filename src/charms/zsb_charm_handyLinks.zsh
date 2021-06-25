#!/usr/bin/env zsh
set -o errexit
set -o nounset
set -o pipefail

declare -r selection=$(print -l "${(@k)linksDictionary}" | rofi -show -dmenu -i -p "Open link")

xdg-open "${linksDictionary[$selection]}"

exit 0

# You must declare `linksDictionary` variable in ~/.zshenv
#
# example
#
# declare -Ar linksDictionary=(
#   [google]='https://www.google.com'
#   [runescape]='https://www.runescape.com'
# )
#
