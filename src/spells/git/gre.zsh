gre() { git restore -- "$@" && ${zsb}.gitStatus; }

hisIgnore gre 'gre .'

compdef "_${zsb}.gitUnrepeat 'unstaged'" gre

