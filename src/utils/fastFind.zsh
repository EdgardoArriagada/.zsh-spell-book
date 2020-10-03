fastFind() {
  local dir="$1"
  local word="$2"
  local type="$3"
  find "$dir" -name "$word" -type "$type" \
      -not -path "*/node_modules/*" \
      -not -path "*/build/*" \
      -not -path "*/dist/*" \
      -not -path "*/public/*" \
      -not -path "*/tmp/*" \
      -not -path "*/\.cache/*" \
      -not -path "*/\.git/*" \
      -print -quit
}
