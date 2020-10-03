# callback function should not receive args
${zsb}_yesNoMenu() {
  local onYesCallback="$1"
  while true; do
    read yn
    case $yn in
    [Yy]*)
      $onYesCallback
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
