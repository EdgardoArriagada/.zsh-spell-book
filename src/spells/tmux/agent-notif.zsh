mf() {
  zsb_tmux_agent_notification --finished _ $TMUX_PANE;
  zsb_play ~/.zsh-spell-book/src/media/sounds/aoe_farm.wav
}

mw() {
  zsb_tmux_agent_notification --working _ $TMUX_PANE;
}

hisIgnore mf mw
