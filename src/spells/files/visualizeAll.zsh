${zsb}.loopFiles() {
  for file in $*; do
    ${zsb}.info "Visualizing $file"
    v "$file"
  done
}

visualizeAll() {
  local extension=$1
  ${zsb}.loopFiles ./*.${extension}
  ${zsb}.tryGitStatus
}

vmatches() {
  local matchingFiles=(`rg --files-with-matches $1`)
  (( $#matchingFiles )) || ${zsb}.throw "No files found matching `hl $1`"
  ${zsb}.loopFiles ${matchingFiles[@]}
  ${zsb}.tryGitStatus
}

vgit() {
  local gitFiles=(`${zsb}.getGitFiles $1`)
  (( $#gitFiles )) || ${zsb}.throw "No files found matching `hl $1`"
  ${zsb}.loopFiles ${gitFiles[@]}
  ${zsb}.tryGitStatus
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


