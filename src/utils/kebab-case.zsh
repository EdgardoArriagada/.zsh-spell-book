kebab-case() {
  local result=$(echo "$1" |
    sed 's/[A-Z]/ &/g;s/^ //' |
    tr -s ' ' |
    sed -e 's/ /-/g' |
    sed -e 's/-$//g')

  echo "${(L)result}"
}
