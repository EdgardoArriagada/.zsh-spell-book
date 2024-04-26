# brew install pandoc

convertMarkdownToOdt() {
  local input=${1:?Error: Input file is required.}
  local output=${input:r}.odt

  [[ ! -f $input ]] && ${zsb}.throw "File $input does not exist"
  [[ "${input:e:l}" != "md" ]] && ${zsb}.throw "File $input does not have .md extension"

  pandoc $input -f markdown -t odt -s -o $output && \
    ${zsb}.success "File `hl $input` converted to `hl $output`"
}

_${zsb}.convertMarkdownToOdt() {
  local mdFiles=( ${(@)$(ls | rg '\.md$')} )

  _describe 'command' mdFiles
}

compdef _${zsb}.convertMarkdownToOdt convertMarkdownToOdt
