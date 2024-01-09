# when $# > 1, $1 is a comma separated list of globs

__dondedice() {
  local maxColumns=`${zsb}.getMaxSearchColumns`

  if [[ -n "$skip" ]];then
     rg --glob "!{package-lock.json,${skip[2]}}" $@ -M $maxColumns
     return $?
  fi

  rg --glob '!{package-lock.json}' $@ -M $maxColumns
}

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
