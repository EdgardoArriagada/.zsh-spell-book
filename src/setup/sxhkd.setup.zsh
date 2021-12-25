if (( $ZSB_MACOS )); then
  ${zsb}.safeLink skhd/.skhdrc ~/.skhdrc
else
  ${zsb}.safeLink sxhkd/sxhkd.desktop ~/.config/autostart/sxhkd.desktop

  ${zsb}.safeLink sxhkd/sxhkdrc ~/.config/sxhkd/sxhkdrc
fi
