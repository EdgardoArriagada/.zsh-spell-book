declare -Ag ZSB_CACHE

${zsb}.cache.set() {
  ZSB_CACHE[${1},VALUE]="$2"
  ZSB_CACHE[${1},START]="$(date +%s)"
  ZSB_CACHE[${1},DURATION]="${3:=5}"
}

${zsb}.cache.get() {
  local -r cacheValue="${ZSB_CACHE[${1},VALUE]}"
  [[ -z "$cacheValue" ]] && return 0

  local -r cacheStart="${ZSB_CACHE[${1},START]}"
  local -r cacheDuration="${ZSB_CACHE[${1},DURATION]}"
  local -r elapsedTime="$(($(date +%s) - $cacheStart))"

  if [[ "$elapsedTime" -le "$cacheDuration" ]]; then
    echo "$cacheValue" && return 0
  fi

  ZSB_CACHE[${1},VALUE]=""
  ZSB_CACHE[${1},START]=""
  ZSB_CACHE[${1},DURATION]=""
}

