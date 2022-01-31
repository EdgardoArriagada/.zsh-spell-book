# vi mode
bindkey -v
export KEYTIMEOUT=1

autoload throw catch

# Partial tab completion color
zstyle -e ':completion:*:default' list-colors 'reply=("${PREFIX:+=(#bi)($PREFIX:t)(?)*==34=34}:${(s.:.)LS_COLORS}")';

## Autosuggestion plugin
bindkey '^ ' autosuggest-accept

## Highlight plugin
typeset -Ag ZSH_HIGHLIGHT_STYLES

# To differentiate aliases and functions from other command types
ZSH_HIGHLIGHT_STYLES[alias]='fg=magenta,bold'
ZSH_HIGHLIGHT_STYLES[function]='fg=blue,bold'

# To have paths colored instead of underlined
ZSH_HIGHLIGHT_STYLES[path]='fg=blue'

# Misc
ZSH_HIGHLIGHT_STYLES[comment]='fg=cyan'

ZSB_CURSOR_DEFAULT='\e[6 q'
ZSB_CURSOR_VIM='\e[1 q'

# Change cursor shape for different vi modes.
function zle-keymap-select {
  if [[ ${KEYMAP} == vicmd ]] ||
     [[ $1 = 'block' ]]; then
    echo -ne "$ZSB_CURSOR_VIM"
  elif [[ ${KEYMAP} == main ]] ||
       [[ ${KEYMAP} == viins ]] ||
       [[ ${KEYMAP} = '' ]] ||
       [[ $1 = 'beam' ]]; then
    echo -ne "$ZSB_CURSOR_DEFAULT"
  fi
}
zle -N zle-keymap-select
zle-line-init() {
    zle vi-insert # initiate `vi insert` as keymap (can be removed if `bindkey -V` has been set elsewhere)
    echo -ne "$ZSB_CURSOR_DEFAULT"
}
zle -N zle-line-init
echo -ne "$ZSB_CURSOR_DEFAULT" # Use beam shape cursor on startup.
preexec() { echo -ne "$ZSB_CURSOR_DEFAULT"; } # Use beam shape cursor for each new prompt.

function append-last-word { ((++CURSOR)); zle insert-last-word; zle vi-insert; }
zle -N append-last-word

# for more, run "zle -al"
bindkey -M vicmd '.' append-last-word
bindkey -M viins '^a' beginning-of-line
bindkey -M viins '^p' history-search-backward
bindkey -M viins '^n' history-search-forward
bindkey -M viins '^e' end-of-line
bindkey -M viins '^[d' kill-word
bindkey -M viins '^w' backward-kill-word
bindkey -M viins '^r' history-incremental-pattern-search-backward
bindkey -M viins '^u' kill-buffer # prevent `Ctrl + u` from not working after entering viins again
bindkey -M viins '^q' push-line

prompt_context() {
  local emojis=( '⚡️' '🔥' '💀' '🦄' '🌈' '🚀' '🌙' '🌎' '🌞' '🪐' )
  local RAND_EMOJI_N=$(( $RANDOM % ${#emojis[@]} + 1))
  prompt_segment black default "${emojis[$RAND_EMOJI_N]}"
}

declare ZSB_HISTORY_IGNORE=(ls)

hisIgnore() ZSB_HISTORY_IGNORE+=( $@ )

# the following line should be applied as an automatic call
# export HISTORY_IGNORE="(${(j:|:)ZSB_HISTORY_IGNORE})"
