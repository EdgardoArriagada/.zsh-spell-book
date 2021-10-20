setWallpaper() {
  [[ ! -d ~/Wallpapers ]] && mkdir ~/Wallpapers
  [[ -z "$(ls ~/Wallpapers)" ]] && ${zsb}.throw "You don't have any wallpapers in $(hl ~/Wallpapers)."

  (cd ~/Wallpapers && ls |
   sxiv -ftio |
   awk 'FNR == 1' |
   xargs -I_ xwallpaper --center ~/Wallpapers/_)
}
