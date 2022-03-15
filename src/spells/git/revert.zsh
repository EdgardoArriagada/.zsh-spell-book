revert() {
  local rev=${1:='HEAD'}

  if [[ "$rev" = 'HEAD' ]]; then
    ${zsb}.warning "This will undo the last commit with another commit. \n
      \rPrefer `hl 'git reset HEAD~'` if commit is `hl not` online."
    ${zsb}.confirmMenu.withPrompt
  fi

  printAndRun "git revert ${rev}"
}

hisIgnore revert

_${zsb}.nocompletion revert
