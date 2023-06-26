${zsb}.stashListMenu() {
  local stashList="$(git stash list)"
  [[ -z "$stashList" ]] && ${zsb}.cancel "Stash list is empty."

  local callback=${1:?'You must provida a callback function'}
  local callbackArg="$2"

  if [[ -n "$callbackArg" ]]; then
    "$callback" "$callbackArg"; return $?
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

