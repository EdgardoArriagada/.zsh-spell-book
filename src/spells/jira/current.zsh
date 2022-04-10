${zsb}.current.getDir() print ${ZSB_TICKETS_DIR}/${ZSB_PARENT_TICKET}

${zsb}_createCurrentDir() {
  ${zsb}.validateJira

  local parentTicketDir=`${zsb}.current.getDir`

  [[ ! -d ${parentTicketDir} ]] && \
    mkdir -p ${parentTicketDir} && \
      (builtin cd ${parentTicketDir} && git init);

  return 0
}

alias cdcurrent="${zsb}_createCurrentDir && cd `${zsb}.current.getDir`"
alias vnotescurrent="${zsb}_createCurrentDir && (cdcurrent && nnvim `${zsb}.current.getDir`/NOTES.md)"
alias cnotescurrent="${zsb}_createCurrentDir && zsb_cat `${zsb}.current.getDir`/NOTES.md"
alias ncurrent="cdcurrent && (( $ZSB_MACOS )) && open . || nautilus ."

hisIgnore cdcurrent vnotescurrent cnotescurrent ncurrent

pomodorocurrent() {
  ${zsb}.validate 'ZSB_CURRENT_TICKET'
  local inputTime="$1"
  shift 1
  ${zsb}.info "Current ticket: `hl ${ZSB_CURRENT_TICKET}`"
  pomodoro "$inputTime" "$ZSB_CURRENT_TICKET" "$*"
}
