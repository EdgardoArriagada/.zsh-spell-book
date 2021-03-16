o() (
  local file="$1"
  local extension=$(${zsb}.extensionOf "$file")

  case "$extension" in
    xlsx|csv|odt)
      libreoffice --calc "$file" ;;
    pdf)
      xdg-open "$file" ;;
    *)
      vim "$file" ;;
  esac
)

