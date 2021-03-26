# extract files
unz() (
  local file="$1"

  case "${file}" in
    *.tar.gz)
      tar -xzf "$file" ;;
    *.gz)
      gunzip --keep "$file" ;;
    *)
      ${zsb}.throw 'Unhandled file extension' ;;
  esac
)

