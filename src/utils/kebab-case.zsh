kebab-case() {
  local result=$(echo "$1" |
    sed 's/[A-Z]/ &/g' | # Add spaces before every capital letter
    sed 's/^ //' | # Remove space at the beginning
    sed 's/[ _]/-/g' | # Replace spaces and underscores by dashes
    tr -s '-' | # Remove duplicated dashes
    sed 's/-$//') # remove final dash if any

  echo "${(L)result}" # Lowercase the result
}
