dondedice() {
  zparseopts -D -E -F -- m=multiCase s:=skip || return 1

  local searchInput=${1:?Error: Search input required.}
  shift

  if [[ -z "$multiCase" ]]; then
    __dondedice $searchInput $@
    return $?
  fi

  local currentKase=`wcase -w "$searchInput" 2>&1`

  if [[ "$currentKase" == 'Invalid input' ]]; then
    ${zsb}.throw "Invalid input for multicase."
  fi

  local -A duplicates

  for kase in $ZSB_WCASE_CASES; do
    local result=`wcase --${kase} -w $searchInput`

    if [[ -n ${duplicates[$result]} ]]
      then continue
      else duplicates[$result]=1
    fi

    __dondedice $result $@
  done

  unset duplicates
}

__dondedice() {
  local maxColumns=`${zsb}.getMaxSearchColumns`
  local globs="-g '!package-lock.json'"

  for glob in ${(z)skip[2]}; do
    globs+=" -g '!$glob'"
  done

  eval "rg $globs '$@' -M $maxColumns"
}

# util that already skip tests and allow to pass multiple other skips
snt() {
  zparseopts -D -E -F -- m=multiCase s:=skip || return 1
  local skips='*__tests__* *Test.java *mocks* *fixtures*'

  for glob in ${(z)skip[2]}; do
    skips+=" $glob"
  done

  dondedice -s $skips $multiCase $@
}

compdef "_${zsb}.nonRepeatedListD \
  '-m:Multicase search' \
  '-s:Skip given space separated list of globs' \
  " dondedice snt
