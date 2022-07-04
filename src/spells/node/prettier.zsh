prettier() {
  zparseopts -D -E -F -- \
    -init=init \
    -install=install \
    -file=file \
    || return 1

  if [[ -n "$init" ]]; then
    _prettier.install
    _prettier.file
    return
  fi

  if [[ -n "$install" ]]; then
    _prettier.install
    return
  fi

  if [[ -n "$file" ]]; then
    _prettier.file
    return
  fi
}

_prettier.file() {
  local -r prettierFile=.prettierrc.yaml
  print "trailingComma: \"es5\"
tabWidth: 2
semi: false
singleQuote: true
" > ${prettierFile} && \
  ${zsb}.success "$(hl ${prettierFile}) initialized" && \
  zsb_cat ${prettierFile}
}

_prettier.install() {
  printAndRun 'npm install --save-dev --save-exact prettier'
}

compdef "_${zsb}.nonRepeatedListD \
  '--init:install and create prettier .prettierrc.yaml' \
  '--install:install prettier' \
  '--file:create .prettierrc.yaml file' \
  " prettier

