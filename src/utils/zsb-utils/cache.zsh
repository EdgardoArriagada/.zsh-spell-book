${zsb}.cache.set() {
  local -r duration=${3:='5'}
  eval "ZSB_CACHE_VALUE__${1}=${2}"
  eval "ZSB_CACHE_START__${1}=$(date +%s)"
  eval "ZSB_CACHE_DURATION__${1}=${duration}"
}

${zsb}.cache.get() {
  local -r cacheValue="$(eval "echo \"\$ZSB_CACHE_VALUE__${1}\"")"
  [[ -z "$cacheValue" ]] && return 0

  local -r cacheStart="$(eval "echo \"\$ZSB_CACHE_START__${1}\"")"
  local -r cacheDuration="$(eval "echo \"\$ZSB_CACHE_DURATION__${1}\"")"

  local -r elapsedTime="$(($(date +%s) - $cacheStart))"
  if [[ "$elapsedTime" -gt "$cacheDuration" ]]; then
    eval "unset ZSB_CACHE_VALUE__${1}"
    eval "unset ZSB_CACHE_START__${1}"
    eval "unset ZSB_CACHE_DURATION__${1}"
    return 0
  fi

  echo "$cacheValue"
}

