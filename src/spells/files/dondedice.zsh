dondedice() {
  zparseopts -D -E -F -- m=multiCase s:=skip

  local searchInput=${1:?Error: Search input required.}
  shift

  if [[ -z "$multiCase" ]]; then
    __dondedice $searchInput $@
    return $?
  fi

  for kase in $ZSB_WCASE_CASES; do
    __dondedice "`wcase --${kase} -w $searchInput`" $@
  done
}

__dondedice() {
  local maxColumns=`${zsb}.getMaxSearchColumns`
  local globs="-g '!package-lock.json'"

  for glob in ${(z)skip[2]}; do
    globs+=" -g '!$glob'"
  done

  eval "rg $globs $@ -M $maxColumns"
  return $?
}

# util that already skip tests and allow to pass multiple other skips
snt() {
  zparseopts -D -E -F -- s:=skip
  local skips='*__tests__* *Test.java'

  for glob in ${(z)skip[2]}; do
    skips+=" $glob"
  done

  dondedice -s $skips $@
}

compdef "_${zsb}.nonRepeatedListD \
  '-m:Multicase search' \
  '-s:Skip given space separated list of globs' \
  " dondedice snt
