createTsReact() {
  local -r projectName=${1:?'You must provide a projectName'}
  npx create-react-app "$projectName" --template typescript && \
    yarn add -D prettier && \
    prettierInit && \
    alert -i 'face-cool' "Happy Hacking" && \
    return 0

  alert -i 'face-sad' 'Something went wrong'
}
