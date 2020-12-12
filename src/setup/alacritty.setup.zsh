# Attemtps to create a symlink from ~/.config/alacritty.yml
# to this projects alacritty.yml

() {
  local LOCAL_ALACRITTY_CONF=${ZSB_DIR}/src/configurations/alacritty-config/alacritty.yml
  local ALACRITTY_CONF=~/.config/alacritty.yml

  if [ "$(readlink -f $ALACRITTY_CONF)" = "$LOCAL_ALACRITTY_CONF" ]; then
    echo "${ZSB_INFO} alacritty setup is already done."
    return 0
  fi

  if [ -f $ALACRITTY_CONF ]; then
    echo "${ZSB_ERROR} $ALACRITTY_CONF is already busy. Please back up it manually before proceeding"
    return 1
  fi

  ln -s ${LOCAL_ALACRITTY_CONF} $ALACRITTY_CONF && \
    echo "${ZSB_SUCCESS} $(hl '$ALACRITTY_CONF') file linked"
}
