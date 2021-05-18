_${zsb}.singleComp() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r comp=( "$@" )
  _describe 'command' comp
}

_${zsb}.singleCompC() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r comp=( $(eval "$1") )
  _describe 'command' comp
}

_${zsb}.cachedSingleCompC() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r cacheKey=${2:?'You must provide a cacheKey'}
  local -r cacheDuration=${3:='5'}
  local comp
  local -r cachedValue="$(${zsb}.cache.get "$cacheKey")"
  if [[ -z "$cachedValue" ]]; then
    comp=( $(eval "$1") )
    ${zsb}.cache.set "$cacheKey" "${comp[*]}" "$cacheDuration"
  else
    comp=( $(echo "$cachedValue") )
  fi
  _describe 'command' comp
}

