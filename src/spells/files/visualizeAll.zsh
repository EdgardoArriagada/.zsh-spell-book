visualizeAll() {
  local -r extension=${1:?'You must provide a file extension'}

  for file in ./*.${extension}; do
    ${zsb}.info "Visualizing ${file}"
    eval "v '${file}'"
  done
}

_${zsb}.visualizeAll() {
  local files=("${(@f)$(ls | rg '\.')}")
  local compList=(${(u)files##*\.})

  _describe 'command' compList
}

compdef _${zsb}.visualizeAll visualizeAll
