o() {
  local -r file="$1"

  case "$file" in
    *.xlsx|*.csv|*.odt)
      libreoffice --calc "$file" ;;
    *.pdf)
      xdg-open "$file" ;;
    *)
      eval "nnvim ${file}" ;;
  esac
}

alias v="o"
