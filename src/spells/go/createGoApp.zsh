createGoApp() {
  local args=${@:?Error: You must provide a project name.}
  local projectName=`wcase -w $args --kebab`

  [[ -d $projectName ]] && ${zsb}.throw "$(hl "$1") is already a dir."

  mkdir $projectName && \
  cd $projectName && \
  git init && \
  go mod init example/user/$projectName && \
  alert -i 'face-cool' "Initialized" && \
  return 0

  alert -i 'face-sad' 'Something went wrong'
}
