${zsb}.confirmMenu() {
  while true; do
    read yn
    case $yn in
    [Yy]*)
      return 0 ;;
    [Nn]*)
      echo "${ZSB_INFO} Cancelled."
      return 1
      ;;
    *)
      echo "${ZSB_PROMPT} Please answer yes or no" ;;
    esac
  done
}
