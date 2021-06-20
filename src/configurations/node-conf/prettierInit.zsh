prettierInit() {
  local -r prettierFile=.prettierrc.yaml
  echo "trailingComma: \"es5\"
tabWidth: 2
semi: false
singleQuote: true
" > ${prettierFile} && \
  ${zsb}.success "$(hl ${prettierFile}) initialized" && \
  batcat ${prettierFile}
}

_${zsb}.nocompletion prettierInit

