parentBranch () {
  git log --first-parent --pretty='%d' \
    | grep 'origin\/' \
    | head -2 \
    | tail -1 \
    | sed 's/,.*//g' | sed 's/origin\///g' | sed 's/[()]//g'
}
