parentBranch () {
  git log --first-parent --pretty='%D' \
    | grep 'origin\/' \
    | sed 's/$/,/g' \
    | sed 's/^tag: .*,//g' \
    | head -2 \
    | tail -1 \
    | sed 's/,.*//g' | sed 's/origin\///g'
}
