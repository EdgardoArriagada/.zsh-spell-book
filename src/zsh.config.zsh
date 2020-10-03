### Enable custom autocompletion
if ! bashcompinit >/dev/null 2>&1; then
  autoload bashcompinit
  bashcompinit
fi
