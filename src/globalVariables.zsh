### GLOBAL VARIABLES ###

# Colors
# usefull link: https://misc.flogisoft.com/bash/tip_colors_and_formatting
ZSB_RED=9
ZSB_GREEN=10
ZSB_YELLOW=11
ZSB_BLUE=12
ZSB_PURPLE=13
ZSB_CYAN=14

# Git
ZSB_GIT_DEFAULT_BRANCHES="(master|develop|devel|dev|main)"
ZSB_GIT_AWARE="--aware:Proceed even if affecting a default branch"
ZSB_GIT_FORCE="--force:Proceed even in commit have already been pushed online"
ZSB_GIT_LOCK_FILES=( package-lock.json Gemfile.lock yarn.lock Cargo.lock go.sum )
ZSB_GIT_PACKAGE_FILES=( package.json Gemfile Cargo.toml go.mod )
declare -gAr ZSB_GIT_FILETYPE_TO_REGEX=(
  ['staged']='^[MARCD]'
  ['green']='^[MARCD]'
  ['unstaged']='^.[MARCD]'
  ['untracked']='^\?{2}'
  ['red-safe']='^.[MARCD\?]'
  ['red']='^.[MARCUD\?]'
  ['red-with-diff']='^.[MARCD]'
  ['unmerged']='(^U)|(^.U)'
)

