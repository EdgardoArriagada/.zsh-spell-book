kebab-case() {
  local result=$(echo "$1" |
    sed 's/[A-Z]/ &/g' | # Add spaces before every capital letter
    sed 's/[ _]/-/g' | # Replace spaces and underscores by dashes
    tr -s '-' | # Remove duplicated dashes
    sed 's/^-//;s/-$//') # Remove dashes at the beginning and at the end

  echo "${(L)result}" # Lowercase the result
}
