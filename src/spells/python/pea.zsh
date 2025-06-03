# poetry env activate
pea() {
  eval 'python --version' # lazy load pyenv
  eval '$(poetry env activate)'
}
