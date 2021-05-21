alias gswd='git checkout develop'
alias gswm='git checkout master'
alias gs='${zsb}.gitStatus'

alias pop='git stash pop'
alias sign='git commit --amend --no-edit --gpg-sign && ${zsb}.gitStatus'
alias glog="git log --color --graph --pretty=format:'%C(magenta)%h%C(reset) -%C(yellow)%d%C(reset) %s %C(green)(%cr) %C(bold blue)<%an>%C(reset)' --abbrev-commit --branches"
alias gsp='git status --porcelain=v2'

alias gst='git status'
alias gf='git fetch'
alias gsh='git show'

alias GSWD='toggleCapsLock && gswd'
alias GSWM='toggleCapsLock && gswm'
alias GS='toggleCapsLock && gs'
alias POP='toggleCapsLock && pop'
alias SIGN='toggleCapsLock && sign'
alias GLOG='toggleCapsLock && glog'
alias GSP='toggleCapsLock && gsp'
alias GST='toggleCapsLock && gst'
alias GF='toggleCapsLock && gf'
alias GSH='toggleCapsLock && gsh'

# this aliases reveice args, so toggleCapsLocks makes no sense
alias gbm='git branch -m'
alias gswa='git switch'
alias gswc='git switch -c'

alias onotif="sendmeto 'https://github.com/notifications'"
