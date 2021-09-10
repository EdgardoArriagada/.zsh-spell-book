__${zsb}.sourceFiles() {
  for file in $*; do
    ${zsb}.fillWithToken '_'
    ${zsb}.info "sourcing $(hl ${file##*/})"
    source "$file"
  done
}

local setupFiles=( ${ZSB_DIR}/src/setup/**/*.setup.zsh )

__${zsb}.sourceFiles $setupFiles

