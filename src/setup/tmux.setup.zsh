# Attemtps to create a symlink from ~/.tmux.conf
# to this projects tmux.conf

() {
  local LOCAL_TMUX_CONF=${ZSB_DIR}/src/configurations/tmux-config/tmux.conf

  if [ "$(readlink -f ~/.tmux.conf)" = "$LOCAL_TMUX_CONF" ]; then
    echo "${ZSB_INFO} tmux setup is already done."
    return 0
  fi

  if [ -f ~/.tmux.conf ]; then
    echo "${ZSB_ERROR} ~/.tmux.conf is already busy. Please back up it manually before proceeding"
    return 1
  fi

  ln -s ${LOCAL_TMUX_CONF} ~/.tmux.conf && \
    echo "${ZSB_SUCCESS} $(hl '~/.tmux.conf') file linked"
}
