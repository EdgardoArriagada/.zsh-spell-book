# TODO Refactor to remove dup code
_${zsb}.singleComp() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r comp=( "$@" )
  _describe 'command' comp
}

# singleComp with Command
_${zsb}.singleCompC() {
  [[ "$CURRENT" -gt 2 ]] && return 0
  local -r comp=( $(eval "$1") )
  _describe 'command' comp
}

# cache expensive commands
_${zsb}.cachedSingleComp() {
  [[ "$CURRENT" -gt 2 ]] && return 0

  local comp
  local -r inputCommand=${1:?'You must provide a command'}
  local -r cacheKey=${2:?'You must provide a cacheKey'}
  local -r cacheDuration=${3:='5'}
  local -r cachedValue="$(${zsb}.cache.get "$cacheKey")"

  if [[ -z "$cachedValue" ]]; then
    comp=( $(eval "$inputCommand") )
    ${zsb}.cache.set "$cacheKey" "${comp[*]}" "$cacheDuration"
  else
    comp=( "${(z)cachedValue}" )
  fi
  _describe 'command' comp
}

# cachedSingleComp by Working Directory
_${zsb}.cachedSingleCompWD() {
  [[ "$CURRENT" -gt 2 ]] && return 0

  local comp
  local -r inputCommand=${1:?'You must provide a command'}
  local -r cacheKey=${2:?'You must provide a cacheKey'}
  local -r cacheDuration=${3:='5'}

  local -r cacheKeyWithWD="${cacheKey}---${PWD}"
  local -r cachedValue="$(${zsb}.cache.get "$cacheKeyWithWD")"

  if [[ -z "$cachedValue" ]]; then
    comp=( $(eval "$inputCommand") )
    ${zsb}.cache.set "$cacheKeyWithWD" "${comp[*]}" "$cacheDuration"
  else
    comp=( "${(z)cachedValue}" )
  fi
  _describe 'command' comp
}

