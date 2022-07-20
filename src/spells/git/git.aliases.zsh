alias gswd='git checkout develop'
alias gs='${zsb}.gitStatus'
alias pop='git stash pop'
alias amend='git commit --amend --no-edit --gpg-sign'
alias glog="git log --color --graph --pretty=format:'%C(magenta)%h%C(reset) -%C(yellow)%d%C(reset) %s %C(green)(%cr) %C(bold blue)<%an>%C(reset)' --abbrev-commit --branches"
alias gsp='git status --porcelain=v2'

hisIgnore gswd gs pop amend glog gsp

alias gst='git status'
alias gf='spinner git fetch'
alias gsh='git show'
alias prv='gh pr view --web'
alias vpr='prv'
alias guu="ggl && gf"

hisIgnore gsw gf gsh prv vpr

alias GSWD='toggleCapsLock && gswd'
alias GS='toggleCapsLock && gs'
alias POP='toggleCapsLock && pop'
alias SIGN='toggleCapsLock && sign'
alias GLOG='toggleCapsLock && glog'
alias GSP='toggleCapsLock && gsp'
alias GST='toggleCapsLock && gst'
alias GF='toggleCapsLock && gf'
alias GSH='toggleCapsLock && gsh'
alias GPV='toggleCapsLock && prv'

# this aliases reveice args, so toggleCapsLocks makes no sense
alias gswa='git switch'
alias gswc='git switch -c'

alias onotif="sendmeto 'https://github.com/notifications'"
