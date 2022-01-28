${zsb}.removeNodeVersionCacheDecorator()
  functions[$1]='
    if [[ "${1} ${2}" = "alias default" ]] && [[ -f ~/temp/cache/node-version ]]; then
      ${zsb}.info "Removing cached node version"
      rm ~/temp/cache/node-version
    fi
    () { '$functions[$1]'; } "$@"
    local ret=$?
    return $ret'

${zsb}.useCachedNodeVersion() {
  mkdir -p ~/temp/cache

  local localVersion=`cat .nvmrc`
  local cachedVersion=`cat ~/temp/cache/node-version`

  if [[ -z "$cachedVersion" ]]; then
    ${zsb}.info "Cached version not found. Local: `hl ${localVersion}`."
  else
    ${zsb}.info "Cached: `hl ${cachedVersion}`. Local: `hl ${localVersion}`."
  fi

  ${zsb}.doesMatch "$cachedVersion" "^${localVersion}.*" && return 0

  eval 'nvm use'

  local currentVersion=`node -v`
  eval "nvm alias default ${currentVersion}"
  printf "$currentVersion" > ~/temp/cache/node-version
}
