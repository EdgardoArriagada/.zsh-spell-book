${zsb}.pomodoro.viewPomodoroLogs() {
  local fileName="$1"
  local fileOpenerProgram="$2"

  # File format: mnt-dd-yyyy.log
  local year="$(echo "$fileName" | cut -d "-" -f3 | cut -d "." -f1)"
  local month="$(echo "$fileName" | cut -d "-" -f 1)"

  local logFolder="${ZSB_DIR}/logs/pomodoro/${year}/${month}"

  local fullFilePath=${logFolder}/${fileName}

  if [[ ! -f $fullFilePath ]]; then
    ${zsb}.info "Pomodoro not found"
    return 0
  fi

  ${fileOpenerProgram} ${fullFilePath}
}

${zsb}.pomodoro.lastWorkingDay() {
  local today=`date +%w`
  local SUNDAY=0
  local MONDAY=1
  if [[ $today = $SUNDAY ]] ; then
    echo "2 days ago"
  elif [[ $today = $MONDAY ]] ; then
    echo "3 days ago"
  else
    echo "1 day ago"
  fi
}


cpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log) zsb_cat
vpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log) nvim

cypomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log -d "yesterday") zsb_cat
vypomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log -d "yesterday") nvim

cwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log -d "$(${zsb}.pomodoro.lastWorkingDay)") zsb_cat
vwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log -d "$(${zsb}.pomodoro.lastWorkingDay)") nvim
