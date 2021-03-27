# extract files
unz() (
  local file="$1"
  [ ! -f "$file" ] && ${zsb}.throw 'Input is not a file.'

  case "${file}" in
    *.tar.gz) tar -xzf "$file" ;;
    *.gz) gunzip --keep "$file" ;;
    *) ${zsb}.throw 'Unhandled file extension.' ;;
  esac
)

