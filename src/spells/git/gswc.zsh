gswc() {
  git switch -c ${1}
}

compdef "_${zsb}.singleCompC 'git branch --show-current'" gswc
