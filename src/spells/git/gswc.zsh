gswc() {
  git switch -c ${@}
}

compdef "_${zsb}.singleCompC 'git branch --show-current'" gswc
