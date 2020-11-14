alias gsw='git switch'
alias gswc='git switch -c'
alias gcd='git checkout develop'
alias gcm='git checkout master'
alias gbm='git branch -m'
alias gsp='git status --porcelain=v2'
alias gs='${zsb}_gitStatus'
alias GS='toggleCapsLock && gs'
alias gst='git status'
alias stash='git stash --include-untracked && ${zsb}_gitStatus'
alias pop='git stash pop'
alias sign='git commit --amend --no-edit --gpg-sign && ${zsb}_gitStatus'
alias glog="git log --color --graph --pretty=format:'%C(magenta)%h%C(reset) -%C(yellow)%d%C(reset) %s %C(green)(%cr) %C(bold blue)<%an>%C(reset)' --abbrev-commit --branches"

alias onotif="sendmeto 'https://github.com/notifications'"

# for the following alias, edgardoArriagada recomend to work in a branch called gh-pages-build
# and remove 'dist' folder from .gitignore in that branch.
# 1.- To build, first move to gh-pages-build branch 		['git checkout gh-pages-build']
# 2.- Merge the branch desired to be build 			['git merge desired-branch']
# 3.- Run you favorite compiler so dist folder is filled	['gulp build'] <- examplle builder, install your own
# 4.- Commit 'git commit' the changes maded to dist folder	['git add .', 'git commit -m "1st build for my-feature"']
# and then run the gh-push command				['gh-push']
alias gh-push='git subtree push --prefix dist origin gh-pages'
