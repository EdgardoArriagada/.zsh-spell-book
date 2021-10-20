setWallpaper() {
  [[ ! -d ~/Wallpapers ]] && mkdir ~/Wallpapers
  [[ -z "$(ls ~/Wallpapers)" ]] && ${zsb}.throw "You don't have any wallpapers in $(hl ~/Wallpapers)."

  local chosenWallpaper=$(builtin cd ~/Wallpapers && ls | sxiv -ftio | awk 'FNR == 1')

  [[ -z "$chosenWallpaper" ]] && return 1
  xwallpaper --center ~/Wallpapers/${chosenWallpaper}

  # Persistence after reboot
  rm -rf ~/temp/wallpaper
  mkdir -p ~/temp/wallpaper
  cp ~/Wallpapers/${chosenWallpaper} ~/temp/wallpaper
}
