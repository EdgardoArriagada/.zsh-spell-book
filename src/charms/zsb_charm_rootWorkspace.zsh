#!/usr/bin/env zsh
set -o errexit
set -o nounset
set -o pipefail

declare -r files=$(ls ${ZSB_CHARM_ROOT_WORKSPACE})

declare dirs
for file in ${(z)files}; do
  [[ -d ${ZSB_CHARM_ROOT_WORKSPACE}/${file} ]] && \
    dirs+="$file\n"
done

declare -r selection=$(echo "$dirs" | rofi -show -dmenu -i -p "Open")

[[ -n "$selection" ]] && \
  code ${ZSB_CHARM_ROOT_WORKSPACE}/${selection}

exit 0

# You must declare `ZSB_CHARM_ROOT_WORKSPACE` variable in ~/.zshenv
#
# example
#
# declare -r ZSB_CHARM_ROOT_WORKSPACE=~/projects
#
