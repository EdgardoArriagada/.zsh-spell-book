parentBranch () {
  git log --first-parent --pretty='%D' \
    | tail -n +2 \
    | rg 'origin/' \
    | sd '$' ',' \
    | sd 'tag: .+?,' '' \
    | head -1 \
    | sd ',.*' '' \
    | sd 'origin/' '' \
    | xargs
}
