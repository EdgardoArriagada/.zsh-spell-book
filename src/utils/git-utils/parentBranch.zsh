parentBranch () {
  git log --pretty=format:'%d' \
    | grep origin/HEAD \
    | sed -e "s/.*origin\/\(.*\), origin\/HEAD.*/\1/"
}
