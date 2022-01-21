autoload throw catch

# Partial tab completion color
zstyle -e ':completion:*:default' list-colors 'reply=("${PREFIX:+=(#bi)($PREFIX:t)(?)*==34=34}:${(s.:.)LS_COLORS}")';

# Use the vi navigation keys in menu completion
bindkey -M menuselect 'h' vi-backward-char
bindkey -M menuselect 'k' vi-up-line-or-history
bindkey -M menuselect 'l' vi-forward-char
bindkey -M menuselect 'j' vi-down-line-or-history

## Autosuggestion plugin
bindkey '^ ' autosuggest-accept

## Highlight plugin
typeset -Ag ZSH_HIGHLIGHT_STYLES

# To differentiate aliases and functions from other command types
ZSH_HIGHLIGHT_STYLES[alias]='fg=magenta,bold'
ZSH_HIGHLIGHT_STYLES[function]='fg=magenta,bold'

# To have paths colored instead of underlined
ZSH_HIGHLIGHT_STYLES[path]='fg=blue'

# Misc
ZSH_HIGHLIGHT_STYLES[comment]='fg=cyan'

prompt_context() {
  local emojis=(⚡️ 🔥 💀 👑 😎 🦄 🌈 🚀 💡 🎉 🌙 🍟 🍔 🌎 ⛄ 🌞 🪐 🔱 🎩 🛸 🎬 🏹 🧜)
  local RAND_EMOJI_N=$(( $RANDOM % ${#emojis[@]} + 1))
  prompt_segment black default "${emojis[$RAND_EMOJI_N]}"
}
