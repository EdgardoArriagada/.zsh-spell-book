watch() {
  local file=$1

  case ${1:e:l} in
    lua) echo $file | entr -c lua $file ;;
    js) node --watch $file ;;
    *) ${zsb}.throw "Not configured yet: `hl $extension`"
  esac
}
