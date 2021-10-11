createTsReact() {
  npx create-react-app "$1" --template typescript && \
    yarn add -D prettier && \
    prettierInit && \
    alert -i 'face-cool' "Happy Hacking" && \
    return 0

  alert -i 'face-sad' 'Something went wrong'
}
