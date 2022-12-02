${zsb}.loopFiles() {
  for file in ${*}; do
    ${zsb}.info "Visualizing ${file}"
    eval "v '${file}'"
  done
}

visualizeAll() {
  local -r extension=${1}
  ${zsb}.loopFiles ./*.${extension}
}

vmatches() {
  local -r matchingFiles=(`rg --files-with-matches ${1}`)
  ${zsb}.loopFiles ${matchingFiles[@]}
}

vgit() {
  local -r gitFiles=(`${zsb}.getGitFiles ${1}`)
  ${zsb}.loopFiles ${gitFiles[@]}
}

_${zsb}.visualizeAll() {
  local files=("${(@f)$(ls | rg '\.')}")
  local uniqueFileExtensions=(${(u)files##*\.})

  _describe 'command' uniqueFileExtensions
}

compdef _${zsb}.visualizeAll visualizeAll

_${zsb}.vgit() {
  local gitFileTypes=(${(@k)ZSB_GIT_FILETYPE_TO_REGEX})
  _describe 'command' gitFileTypes
}

compdef _${zsb}.vgit vgit


