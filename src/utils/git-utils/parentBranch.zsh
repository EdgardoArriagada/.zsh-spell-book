parentBranch () {
  git log --first-parent --pretty='%D' \
    | awk '(NR>1)' \
    | grep 'origin\/' \
    | sed 's/$/,/g' \
    | perl -pe 's|tag: .+?,||g' \
    | head -1 \
    | sed 's/,.*//g' \
    | sed 's/origin\///g' \
    | xargs
}
