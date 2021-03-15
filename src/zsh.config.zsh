# Enable custom autocompletion
if ! bashcompinit >/dev/null 2>&1; then
  autoload bashcompinit
  bashcompinit
fi

bindkey '^e' edit-command-line

# Partial tab completion color
zstyle -e ':completion:*:default' list-colors 'reply=("${PREFIX:+=(#bi)($PREFIX:t)(?)*==34=34}:${(s.:.)LS_COLORS}")';

# Use the vi navigation keys in menu completion
bindkey -M menuselect 'h' vi-backward-char
bindkey -M menuselect 'k' vi-up-line-or-history
bindkey -M menuselect 'l' vi-forward-char
bindkey -M menuselect 'j' vi-down-line-or-history

