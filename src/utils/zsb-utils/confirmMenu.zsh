${zsb}.confirmMenu() {
  while true; do
    read yn

    case $yn in
      [Yy]*) return 0 ;;
      [Nn]*) ${zsb}.cancel "Cancelled." ;;
      *) ${zsb}.prompt "Please answer yes or no" ;;
    esac
  done
}

${zsb}.confirmMenu.continued() {
  while true; do
    read yn

    case $yn in
      [Yy]*) return 0 ;;
      [Nn]*) return 1 ;;
      *) ${zsb}.prompt "Please answer yes or no" ;;
    esac
  done
}

${zsb}.confirmMenu.withItems() {
  local inputWarning="$1"
  local inputItems="$2"

  ${zsb}.warning "$inputWarning"
  echo " "
  [[ -n "$inputItems" ]] && echo "$(hl ${inputItems})"
  echo " "

  ${zsb}.confirmMenu.withPrompt
}

${zsb}.confirmMenu.withPrompt() {
  ${zsb}.prompt "Are you sure? $(hl "[Y/n]")"
  ${zsb}.confirmMenu
}

${zsb}.confirmMenu.warning() {
  ${zsb}.warning "$@"

  ${zsb}.confirmMenu.withPrompt
}

