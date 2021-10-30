${zsb}.stashListMenu() {
  local -r stashList="$(git stash list)"
  [[ -z "$stashList" ]] && ${zsb}.cancel "Stash list is empty."

  local -r callback=${1:?'You must provida a callback function'}
  local -r callbackArg="$2"

  if [[ -n "$callbackArg" ]]; then
    "$callback" "$callbackArg"
    return 0
  fi

  gsl
  ${zsb}.prompt "Enter stash $(hl number)."
  local choosenStashNumber
  read choosenStashNumber
  "$callback" "$choosenStashNumber"
}

${zsb}.stashListMenu.show() git stash show -p "stash@{${1}}"
${zsb}.stashListMenu.apply() git stash apply "$1"

_${zsb}.nocompletion ${zsb}.stashListMenu

alias show="${zsb}.stashListMenu '${zsb}.stashListMenu.show'"
alias apply="${zsb}.stashListMenu '${zsb}.stashListMenu.apply'"

