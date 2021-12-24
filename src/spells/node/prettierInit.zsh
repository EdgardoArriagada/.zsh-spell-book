prettierInit() {
  local -r prettierFile=.prettierrc.yaml
  echo "trailingComma: \"es5\"
tabWidth: 2
semi: false
singleQuote: true
" > ${prettierFile} && \
  ${zsb}.success "$(hl ${prettierFile}) initialized" && \
  zsb_cat ${prettierFile}
}

_${zsb}.nocompletion prettierInit

