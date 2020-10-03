# Url encode source https://www.w3schools.com/tags/ref_urlencode.asp
# Limitations: cannot double quote
googleSearch() (
  declare -A urlEncMap
  urlEncMap[#]="%23"
  urlEncMap[$]="%24"
  urlEncMap[%]="%25"
  local aux="&"
  urlEncMap[$aux]="%26"
  local aux="'"
  urlEncMap[$aux]="%27"
  urlEncMap[+]="%2B"
  urlEncMap[,]="%2C"
  urlEncMap[/]="%2F"
  urlEncMap[:]="%3A"
  local aux=";"
  urlEncMap[$aux]="%3B"
  urlEncMap[=]="%3D"
  urlEncMap[?]="%3F"
  local aux="@"
  urlEncMap[$aux]="%40"
  local aux="["
  urlEncMap[$aux]="%5B"
  local aux="\\"
  urlEncMap[$aux]="%5C"
  local aux="]"
  urlEncMap[$aux]="%5D"
  urlEncMap[^]="%5E"
  local aux="\`"
  urlEncMap[$aux]="%21"
  local aux="{"
  urlEncMap[$aux]="%7B"
  local aux="|"
  urlEncMap[$aux]="%7C"
  local aux="{"
  urlEncMap[$aux]="%7D"

  local input="${*}"
  local searchString=""

  main () {
    sanatizeInput
    createSearchString
    performGoogleSearch
  }

  sanatizeInput() {
    local tmp="$input" # The loop will consume the variable, so make a temp copy first
    local sanatizedString=""
    while [ -n "$tmp" ]; do
      rest="${tmp#?}" # All but the first character of the string
      firstChar="${tmp%"$rest"}" # Remove $rest, and you're left with the first character
      local sanatizedChar="${urlEncMap[$firstChar]}"
      [ -z "$sanatizedChar" ] && sanatizedChar="$firstChar"
      sanatizedString="${sanatizedString}${sanatizedChar}"
      tmp="$rest"
    done
    input="$sanatizedString"
  }

  createSearchString() searchString=$(echo "$input" | sed "s/ /\+/g" | sed -E "s/\++$//g")

  performGoogleSearch() {
    local googleSearchUrl="https://www.google.com/search?q="
    sendmeto "${googleSearchUrl}${searchString}"
  }

  main "$@"
)

alias e="googleSearch"
