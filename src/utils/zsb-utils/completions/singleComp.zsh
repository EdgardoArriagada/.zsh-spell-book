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
  local -r isCacheAlive="$(${zsb}.cache.get "$cacheKey")"
  if [[ "$isCacheAlive" == true ]]; then
    comp=( "${ZSB_CACHED_COMPLETION[@]}" )
  else
    comp=( $(eval "$1") )
    ZSB_CACHED_COMPLETION=( "${comp[@]}" )
    ${zsb}.cache.set "$cacheKey" true "$cacheDuration"
  fi
  _describe 'command' comp
}

