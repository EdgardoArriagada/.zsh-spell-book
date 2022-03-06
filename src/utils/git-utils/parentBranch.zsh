parentBranch () {
  git log --first-parent --pretty='%d' \
    | grep 'origin\/' \
    | sed 's/(tag: .*)//g' \
    | sed '/^[[:space:]]*$/d' \
    | head -2 \
    | tail -1 \
    | sed 's/[()]//g' | sed 's/,.*//g' | sed 's/origin\///g'
}
