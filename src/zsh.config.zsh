# Editor for `fc` command
FCEDIT=nvim

# Disable suspend the output "feature"
stty stop undef # disables C-s
stty start '^-' stop '^-' # disables C-q

setopt INC_APPEND_HISTORY # Write to history as soon as input gets entered
setopt INTERACTIVE_COMMENTS # Enable comments in interactive shells

# Make completion smarter when pressing tab based on current input
zstyle ':completion:*' matcher-list \
    'm:{[:lower:]}={[:upper:]}' \
    '+r:|[._-]=* r:|=*' \
    '+l:|=*'

# prevent unknown commands to be output to the history
zshaddhistory() { whence ${${(z)1}[1]} >| /dev/null || return 1 }

# vi mode
bindkey -v
export KEYTIMEOUT=1

# completion
fpath=(~/.zsh/completions $fpath)

autoload -Uz compinit
compinit

autoload -U bashcompinit
bashcompinit

# allow to use throw
autoload throw catch

## Autosuggestion plugin
bindkey '^ ' autosuggest-accept
ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=#928374"

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
ZSB_CURSOR_VIM='\e[2 q'

# Change cursor shape for different vi modes.
function zle-keymap-select { [[ "$KEYMAP" = "vicmd" || "$1" = "block" ]] && printf $ZSB_CURSOR_VIM || printf $ZSB_CURSOR_DEFAULT; }

zle -N zle-keymap-select
zle-line-init() {
    zle vi-insert # initiate `vi insert` as keymap (can be removed if `bindkey -V` has been set elsewhere)
    printf $ZSB_CURSOR_DEFAULT
}
zle -N zle-line-init
printf $ZSB_CURSOR_DEFAULT # Print startup cursor
preexec() { printf $ZSB_CURSOR_DEFAULT; } # Print cursor for each new prompt

function append-last-word { ((++CURSOR)); zle insert-last-word; zle vi-insert; }
zle -N append-last-word

vim-like-delete-word() {
    local WORDCHARS="${WORDCHARS//[\/:.]/}"

    zle backward-delete-word
}
zle -N vim-like-delete-word

# for more, run "zle -al" and/or "bindkey -l"
bindkey -M vicmd '.' append-last-word

bindkey -M viins '^[c' vi-cmd-mode
bindkey -M viins '^a' beginning-of-line
bindkey -M viins '^p' history-search-backward
bindkey -M viins '^n' history-search-forward
bindkey -M viins '^e' end-of-line
bindkey -M viins '^w' vim-like-delete-word
bindkey -M viins '^v' backward-kill-word
bindkey -M viins '^r' history-incremental-pattern-search-backward
bindkey -M viins '^q' push-line
bindkey -M viins '^o' clear-screen

# Workarouds
bindkey -M viins '^u' kill-buffer # prevent `Ctrl + u` from not working after entering viins again
bindkey "^?" backward-delete-char # makes backspace work as normal when reentering to ins mode in the middle of the command

# Initialize hisignore with some stuff to ignore already
declare ZSB_HISTORY_IGNORE=(
  'l[a,l,s,h,]*'
  'neofetch'
  ' *'
  'l'
)

hisIgnore() ZSB_HISTORY_IGNORE+=( $@ )

# the following line should be applied as an automatic call
# export HISTORY_IGNORE="(${(j:|:)ZSB_HISTORY_IGNORE})"
#
# Alien prompt https://github.com/eendroroy/alien
#
source ~/alien/alien.zsh

export ALIEN_SECTIONS_LEFT=(
  exit
  battery
  path
  vcs_branch:async
  vcs_status:async
  vcs_dirty:async
  newline
  ssh
  venv
  prompt
)

export ALIEN_THEME="bnw"

export ALIEN_SECTION_VCS_DIRTY_FG=0
export ALIEN_SECTION_VCS_DIRTY_BG=4
export ALIEN_GIT_TRACKED_COLOR=0
export ALIEN_GIT_UN_TRACKED_COLOR=0

export ALIEN_GIT_ADD_SYM=⭑⭑
export ALIEN_GIT_DEL_SYM=⭑⭑
export ALIEN_GIT_MOD_SYM=⭑⭑
export ALIEN_GIT_NEW_SYM=⭑⭑
export ALIEN_GIT_PUSH_SYM=' ▲'
export ALIEN_GIT_SYM=''
export ALIEN_GIT_STASH_SYM='⧈ '
