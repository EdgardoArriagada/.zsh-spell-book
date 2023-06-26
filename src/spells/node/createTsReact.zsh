createTsReact() {
  local args="${@:?'You must provide a project name'}"
  local projectName="$(kebab-case "$args")"

  [[ -d $projectName ]] && ${zsb}.throw "$(hl "$1") is already a dir."
  eval "npx create-react-app ${projectName} --template typescript" && \
    cd ${projectName} && \
    prettierHelper --init && \
    alert -i 'face-cool' "Happy Hacking" && \
    return 0

  alert -i 'face-sad' 'Something went wrong'
}
