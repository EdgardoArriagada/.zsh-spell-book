# Attemtps to create a symlink from ~/.tmux.conf
# to this projects tmux.conf

() {
  local LOCAL_TMUX_CONF=${ZSB_DIR}/src/configurations/tmux-config/tmux.conf

  # if: is this config already set, return 0
  [ "$(readlink -f ~/.tmux.conf)" = "$LOCAL_TMUX_CONF" ] && return 0

  if [ -f ~/.tmux.conf ]; then
    echo "${ZSB_INFO} ~/.tmux.conf is already busy."
    return 0
  fi

  ln -s ${LOCAL_TMUX_CONF} ~/.tmux.conf && \
    echo "${ZSB_SUCCESS} $(hl '~/.tmux.conf') file linked"
}
