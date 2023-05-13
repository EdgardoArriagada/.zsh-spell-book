kebab-case() {
  local result=`
    sed 's/[A-Z]/ &/g' <<< $1 | # Add spaces before every capital letter
    sed 's/[ _]/-/g' | # Replace spaces and underscores by dashes
    tr -s '-' | # Remove duplicated dashes
    sed 's/^-//;s/-$//'` # Remove dashes at the beginning and at the end

  print ${(L)result}
}
