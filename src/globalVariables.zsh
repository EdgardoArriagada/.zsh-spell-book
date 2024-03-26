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
declare ZSB_GIT_MAIN_BRANCHES_ARRAY=(master main)
declare ZSB_GIT_DEVELOP_BRANCHES_ARRAY=(develop)

ZSB_GIT_AWARE="--aware:Proceed even if affecting a default branch"
ZSB_GIT_FORCE="--force:Proceed even in commit have already been pushed online"

ZSB_GIT_LOCK_FILES=( package-lock.json pnpm-lock.yaml Gemfile.lock yarn.lock Cargo.lock go.sum manifest.toml )
ZSB_GIT_PACKAGE_FILES=( package.json Gemfile Cargo.toml go.mod gleam.toml )

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

# wcase -h
declare -g ZSB_WCASE_CASES=(
  'flat' 'upper' 'camel' 'pascal' 'snake' 'all-caps' 'kebab' 'train' 'spaced' 'http-header' 'title'
)

### DONT TOUCH ###
declare ZSB_GIT_DEFAULT_BRANCHES_ARRAY=( ${ZSB_GIT_MAIN_BRANCHES_ARRAY[@]} )
ZSB_GIT_DEFAULT_BRANCHES_ARRAY+=( ${ZSB_GIT_DEVELOP_BRANCHES_ARRAY[@]} )

ZSB_GIT_MAIN_BRANCHES="(${(j:|:)ZSB_GIT_MAIN_BRANCHES_ARRAY})"
ZSB_GIT_DEVELOP_BRANCHES="(${(j:|:)ZSB_GIT_DEVELOP_BRANCHES_ARRAY})"
ZSB_GIT_DEFAULT_BRANCHES="(${(j:|:)ZSB_GIT_DEFAULT_BRANCHES_ARRAY})"
