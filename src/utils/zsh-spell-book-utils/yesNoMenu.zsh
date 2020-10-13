${zsb}_yesNoMenu() {
  local onYesCallback="$@"
  while true; do
    read yn
    case $yn in
    [Yy]*)
      [ -z "$onYesCallback" ] && return 0

      eval "$onYesCallback"
      return $?
      ;;
    [Nn]*)
      echo "${ZSB_INFO} Cancelled"
      return 1
      ;;
    *)
      echo "${ZSB_PROMPT} Please answer yes or no"
      ;;
    esac
  done
}
