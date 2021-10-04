${zsb}.confirmMenu() {
  while true; do
    read yn

    case $yn in
      [Yy]*) return 0 ;;
      [Nn]*) ${zsb}.info "Cancelled."; return 1 ;;
      *) ${zsb}.prompt "Please answer yes or no" ;;
    esac
  done
}
