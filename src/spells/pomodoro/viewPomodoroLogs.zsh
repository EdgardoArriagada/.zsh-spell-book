${zsb}_viewPomodoroLogs() {
  local fileName="$1"
  local fileOpenerProgram="$2"

  # File format: mnt-dd-yyyy.log
  local year="$(echo "$fileName" | cut -d "-" -f3 | cut -d "." -f1)"
  local month="$(echo "$fileName" | cut -d "-" -f 1)"

  local logFolder="${ZSB_DIR}/logs/pomodoro/${year}/${month}"

  local fullFilePath=${logFolder}/${fileName}

  if [ ! -f $fullFilePath ]; then
    echo "${ZSB_INFO} Pomodoro not found"
    return 0
  fi

  ${fileOpenerProgram} ${fullFilePath}
}

${zsb}_lastWorkingDay() {
  local today=`date +%w`
  local SUNDAY=0
  local MONDAY=1
  if [ $today = $SUNDAY ] ; then
    echo "2 days ago"
  elif [ $today = $MONDAY ] ; then
    echo "3 days ago"
  else
    echo "1 day ago"
  fi
}


cpomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log) cat
vpomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log) vim

cypomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log -d "yesterday") cat
vypomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log -d "yesterday") vim

cwpomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log -d "$(${zsb}_lastWorkingDay)") cat
vwpomodoto() ${zsb}_viewPomodoroLogs $(date +%b-%d-%Y.log -d "$(${zsb}_lastWorkingDay)") vim
