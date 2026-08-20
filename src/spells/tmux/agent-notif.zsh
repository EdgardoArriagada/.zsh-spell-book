mf() { zsb_tmux_agent_notification --finished _ $TMUX_PANE; }

mw() { zsb_tmux_agent_notification --working _ $TMUX_PANE; }

hisIgnore mf mw
