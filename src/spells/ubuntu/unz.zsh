# extract files
unz() (
  local file="$1"
  [[ ! -f "$file" ]] && ${zsb}.throw 'Input is not a file.'

  case "${file}" in
    *.tar.gz) tar -xzf "$file" ;;
    *.tar.bz2) tar -xjf $file ;;
    *.tar.xz) tar zxvf "$file";;
    *.bz2) bunzip2 "$file" ;;
    *.rar) rar x "$file" ;;
    *.gz) gunzip --keep "$file" ;;
    *.tar) tar xf "$file";;
    *.tbz2) ar xjf "$file" ;;
    *.tgz) tar xzf "$file" ;;
    *.xz) xz -d "$file" ;;
    *.zip) unzip "$file" ;;
    *.Z) uncompress "$file" ;;
    *.7z) 7z x "$file" ;;
    *) ${zsb}.throw 'Unhandled file extension.' ;;
  esac
)

