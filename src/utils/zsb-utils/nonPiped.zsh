# decide if certain value should be taken based on if command is beign piped Ie: "cmd | myCmd"
${zsb}.nonPiped() {
  # Overloading
  case "$#" in
    1) # If piped, take pipe value, otherwise input arg
      [[ ! -p /dev/stdin ]] && echo "$1" || echo "$(< /dev/stdin)" ;;
    2) # If piped, take position of 1st arg (because 1sr arg is now piped instead), otherwise 2d
      [[ ! -p /dev/stdin ]] && echo "$1" || echo "$2" ;;
  esac
}
