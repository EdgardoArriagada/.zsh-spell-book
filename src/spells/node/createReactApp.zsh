__createReactApp() {
  local template=$1
  shift
  local args=${@:?Error: You must provide a project name.}
  local projectName=`wcase -w $args --kebab`

  [[ -d $projectName ]] && ${zsb}.throw "$(hl "$1") is already a dir."

  npm create vite@latest $projectName -- --template $template && \
    cd $projectName && \
    prettierHelper --init && \
    alert -i 'face-cool' "Happy Hacking" && \
    return 0

  alert -i 'face-sad' 'Something went wrong'
}

alias createTsReactApp='__createReactApp ts-react'
alias createReactApp='__createReactApp react'
