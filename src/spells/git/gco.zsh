# I only use this for checking out files
gco() {
  git checkout -- "$@" && ${zsb}.gitStatus
}

hisIgnore gco 'gco .'

compdef "_${zsb}.gitUnrepeat 'unstaged'" gco

