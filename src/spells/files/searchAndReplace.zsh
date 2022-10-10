searchAndReplace() {
  local search=$1
  local replace=$2

  sd ${search} ${replace} `rg --files-with-matches ${search}`
}


