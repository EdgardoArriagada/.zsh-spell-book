${zsb}.current.getDir() print ${ZSB_TICKETS_DIR}/${ZSB_CURRENT_TICKET}

${zsb}.current.validate() {
  ${zsb}.validate 'ZSB_TICKETS_DIR'
  ${zsb}.validate 'ZSB_CURRENT_TICKET'
}

${zsb}_createCurrentDir() {
  ${zsb}.current.validate

  local current_ticket_dir=`${zsb}.current.getDir`

  [[ ! -d ${current_ticket_dir} ]] && \
    mkdir -p ${current_ticket_dir} && \
      (builtin cd ${current_ticket_dir} && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cd `${zsb}.current.getDir`"
alias vnotescurrent="${zsb}_createCurrentDir && (cdcurrent && nnvim `${zsb}.current.getDir`/NOTES.md)"
alias cnotescurrent="${zsb}_createCurrentDir && zsb_cat `${zsb}.current.getDir`/NOTES.md"
alias ncurrent="cdcurrent && (( $ZSB_MACOS )) && open . || nautilus ."

pomodorocurrent() {
  ${zsb}.validate 'ZSB_CURRENT_TICKET'
  local inputTime="$1"
  shift 1
  ${zsb}.info "Current ticket: `hl ${ZSB_CURRENT_TICKET}`"
  pomodoro "$inputTime" "$ZSB_CURRENT_TICKET" "$*"
}
