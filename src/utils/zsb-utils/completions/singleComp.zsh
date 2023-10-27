# TODO Refactor to remove dup code
_${zsb}.singleComp() {
  (( $CURRENT > 2 )) && return 0
  local comp=( "$@" )
  _describe 'command' comp
}

# singleComp with Command
_${zsb}.singleCompC() {
  (( $CURRENT > 2 )) && return 0
  local comp=( $(eval "$1") )
  _describe 'command' comp
}

# cache expensive commands
_${zsb}.cachedSingleComp() {
  (( $CURRENT > 2 )) && return 0

  local comp
  local inputCommand=${1?Error: You must provide a command.}
  local cacheKey=${2?Error: You must provide a cacheKey.}
  local cacheDuration=${3:='5'}
  local cachedValue="$(${zsb}.cache.get "$cacheKey")"

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
  (( $CURRENT > 2 )) && return 0

  local comp
  local inputCommand=${1?Error: You must provide a command.}
  local cacheKey=${2?Error: You must provide a cacheKey.}
  local cacheDuration=${3:='5'}

  local cacheKeyWithWD="${cacheKey}---${PWD}"
  local cachedValue="$(${zsb}.cache.get "$cacheKeyWithWD")"

  if [[ -z "$cachedValue" ]]; then
    comp=( $(eval "$inputCommand") )
    ${zsb}.cache.set "$cacheKeyWithWD" "${comp[*]}" "$cacheDuration"
  else
    comp=( "${(z)cachedValue}" )
  fi
  _describe 'command' comp
}

