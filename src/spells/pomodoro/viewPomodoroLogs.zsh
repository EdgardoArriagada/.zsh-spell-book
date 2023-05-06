${zsb}.pomodoro.viewPomodoroLogs() {
  local fileName=$1
  local fileOpenerProgram=$2

  # File format: mnt-dd-yyyy.log
  local year=`cut -d '-' -f3 <<< ${fileName} | cut -d '.' -f1`
  local month=`cut -d '-' -f 1 <<< ${fileName}`

  local logFolder=${ZSB_DIR}/logs/pomodoro/${year}/${month}

  local fullFilePath=${logFolder}/${fileName}

  if [[ ! -f ${fullFilePath} ]]
    then ${zsb}.cancel 'Pomodoro not found'
  fi

  ${fileOpenerProgram} ${fullFilePath}
}

${zsb}.pomodoro.twoDaysAgo() {
  (( $ZSB_MACOS )) && printf '-v -2d' || printf  '-d 2 days ago'
}

${zsb}.pomodoro.threeDaysAgo() {
  (( $ZSB_MACOS )) && printf '-v -3d' || printf  '-d 3 days ago'
}

${zsb}.pomodoro.oneDayAgo() {
  (( $ZSB_MACOS )) && printf '-v -1d' || printf  '-d 1 day ago'
}

${zsb}.pomodoro.lastWorkingDay() {
  local today=`date +%w`
  local SUNDAY=0
  local MONDAY=1

  if (( ${today} == ${SUNDAY} ))
    then ${zsb}.pomodoro.twoDaysAgo

  elif (( ${today} == ${MONDAY} ))
    then ${zsb}.pomodoro.threeDaysAgo

  else ${zsb}.pomodoro.oneDayAgo

  fi
}

cpomodoro() ${zsb}.pomodoro.viewPomodoroLogs `date +%b-%d-%Y.log` zsb_cat
vpomodoro() ${zsb}.pomodoro.viewPomodoroLogs `date +%b-%d-%Y.log` nvim

cypomodoro() ${zsb}.pomodoro.viewPomodoroLogs `date +%b-%d-%Y.log -d yesterday` zsb_cat
vypomodoro() ${zsb}.pomodoro.viewPomodoroLogs `date +%b-%d-%Y.log -d yesterday` nvim

if (( $ZSB_MACOS )); then
  cwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date `${zsb}.pomodoro.lastWorkingDay` +%b-%d-%Y.log) zsb_cat
  vwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date `${zsb}.pomodoro.lastWorkingDay` +%b-%d-%Y.log) nvim
else
  cwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log `${zsb}.pomodoro.lastWorkingDay`) zsb_cat
  vwpomodoro() ${zsb}.pomodoro.viewPomodoroLogs $(date +%b-%d-%Y.log `${zsb}.pomodoro.lastWorkingDay`) nvim
fi
