declare -gAr ZSH_VERSION_DIC=(
  [tmux]='tmux -V'
  [choose]='choose -v'
  [rust]='hasRust'
  [go]='hasGo'
)

version() {
  local input=${1?Error: No input given.}
  local cmd=${ZSH_VERSION_DIC[$input]}

  if [[ -n "$cmd" ]]
    then printAndRun "$cmd"
    else printAndRun "$input --version"
  fi
}

alias vs=version

compdef "_${zsb}.singleCompC 'compgen -c'" version

