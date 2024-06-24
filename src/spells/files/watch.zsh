watch() {
  local file=${1:?Error: you must provide a file to watch.}

  [[ ! -a $file ]] && ${zsb}.throw "File not found: `hl $file`"

  case ${1:e:l} in
    lua) echo $file | entr -c lua $file ;;
    js) node --watch $file ;;
    *) ${zsb}.throw "Not configured yet: `hl $extension`"
  esac
}
