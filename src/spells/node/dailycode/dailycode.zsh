dailycode() {
  if [ -z "$1" ]; then
    echo "${ZSB_PROMPT} Enter the new repository's name:"
    read problemName
  else
    problemName="$1"
  fi

  cp -pr ${ZSB_DIR}/src/spells/node/dailycode/template/ ./${problemName} &&
    cd ${problemName} && git init && eval "npm install jest"
}
