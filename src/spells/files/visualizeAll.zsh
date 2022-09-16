visualizeAll() {
  local -r extension=${1:?'You must provide a file extension'}

  for file in ./*.${extension}; do
    ${zsb}.info "Visualizing ${file}"
    eval "v '${file}'"
  done
}

_${zsb}.visualizeAll() {
  local files=("${(@f)$(ls | rg '\.')}")
  local uniqueFileExtensions=(${(u)files##*\.})

  _describe 'command' uniqueFileExtensions
}

compdef _${zsb}.visualizeAll visualizeAll
