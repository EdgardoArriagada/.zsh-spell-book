parentBranch () {
  git log --first-parent --pretty='%D' \
    | awk '(NR>1)' \
    | rg 'origin/' \
    | sd '$' ',' \
    | sd 'tag: .+?,' '' \
    | head -1 \
    | sd ',.*' '' \
    | sd 'origin/' '' \
    | xargs
}
