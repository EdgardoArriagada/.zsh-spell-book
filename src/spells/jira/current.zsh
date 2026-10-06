${zsb}_createCurrentDir() {
  ${zsb}.assertJira

  [[ -d $ZSB_CURRENT_TICKET_DIR ]] && return 0

  mkdir -p $ZSB_CURRENT_TICKET_DIR
  print "# $ZSB_CURRENT_LABEL" > $ZSB_CURRENT_TICKET_DIR/NOTES.md
  (builtin cd $ZSB_CURRENT_TICKET_DIR && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cds $ZSB_CURRENT_TICKET_DIR"
alias vcurrent="nvim '+cd $ZSB_CURRENT_TICKET_DIR' $ZSB_CURRENT_TICKET_DIR/NOTES.md"
alias cdc='cdcurrent'

alias vanchors="nvim $ZSB_CURRENT_TICKET_DIR/.anchors"
alias anchor="${zsb}.anchor"

${zsb}.anchor() {
  local anchors=$ZSB_CURRENT_TICKET_DIR/.anchors
  local cwd=$(print -P %~)

  if [[ -f "$anchors" ]] && grep -Fqx -- "$cwd" "$anchors"; then
    ${zsb}.info 'Already in anchors; skipped.'
    return 0
  fi

  print -r -- "$cwd" >> "$anchors" && ${zsb}.success 'Added to anchors!'
}

cda() {
  local anchors=$ZSB_CURRENT_TICKET_DIR/.anchors
  [[ ! -f $anchors ]] && ${zsb}.throw "`hl $anchors` not set"

  local cwd=$(print -P %~)
  local -a entries=("${(@f)$(grep -v '^[[:space:]]*$' "$anchors" | grep -vFx -- "$cwd")}")
  local -a choices
  local i dir repo branch label
  for (( i = 1; i <= ${#entries}; i++ )); do
    dir=${entries[i]/#\~/$HOME}
    label=${entries[i]}
    if repo=$(builtin cd -- "$dir" 2>/dev/null && get_repo_name 2>/dev/null) && [[ -n "$repo" ]]; then
      branch=$(git -C "$dir" branch --show-current 2>/dev/null)
      [[ -n "$branch" ]] || branch=$(git -C "$dir" rev-parse --short HEAD 2>/dev/null)
      label="[$repo - $branch]"
      [[ $(git -C "$dir" rev-parse --git-dir 2>/dev/null) != $(git -C "$dir" rev-parse --git-common-dir 2>/dev/null) ]] && label+=' 󰙅'
    fi
    choices+=( "$i"$'\t'"$label" )
  done
  local selection=$(print -rl -- "${choices[@]}" | fzf --select-1 --delimiter=$'\t' --with-nth=2..)
  [[ -n "$selection" ]] && selection=${entries[${selection%%$'\t'*}]}
  selection=${selection/#\~/$HOME}

  if [[ -z "$selection" ]]
    then ${zsb}.cancel 'You did not move.'
  elif [[ -d "$selection" ]]
    then cd -- "$selection"
  else
    ${zsb}.throw "`hl $selection` is not a directory"
  fi
}

grl() {
  local revLine="`get_repo_name 2>/dev/null` - `git branch --show-current 2>/dev/null`"
  zsb_clipcopy <<< $revLine
  ${zsb}.success "`hl $revLine` copied!"
}

hisIgnore cdcurrent vnotescurrent cnotescurrent ncurrent cinit cdc cda vanchors anchor get-rev-line grl

pomodorocurrent() {
  ${zsb}.assertIsSet 'ZSB_CURRENT_TICKET'
  local inputTime=${1?Error: missing time}

  ${zsb}.info "`hl $ZSB_CURRENT_TICKET` $ZSB_CURRENT_LABEL"
  pomodoro $inputTime $ZSB_CURRENT_TICKET $ZSB_CURRENT_LABEL
}
