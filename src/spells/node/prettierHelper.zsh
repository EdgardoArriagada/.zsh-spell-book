prettierHelper() {
  zparseopts -D -E -F -- \
    -init=init \
    -install=install \
    -file=file \
    || return 1

  if [[ -n "$init" ]]; then
    _prettierHelper.install
    _prettierHelper.file
    return
  fi

  if [[ -n "$install" ]]; then
    _prettierHelper.install
    return
  fi

  if [[ -n "$file" ]]; then
    _prettierHelper.file
    return
  fi

  ${zsb}.throw "you have to pass either --init, --install or --file"
}

_prettierHelper.file() {
  local prettierFile=.prettierrc.yaml
  print "trailingComma: \"es5\"
tabWidth: 2
semi: false
singleQuote: true
" > ${prettierFile} && \
  ${zsb}.success "$(hl ${prettierFile}) initialized" && \
  zsb_cat ${prettierFile}
}

_prettierHelper.install() {
  printAndRun 'npm install --save-dev --save-exact prettier'
}

compdef "_${zsb}.nonRepeatedListD \
  '--init:install and create prettier .prettierrc.yaml' \
  '--install:install prettier' \
  '--file:create .prettierrc.yaml file' \
  " prettierHelper

