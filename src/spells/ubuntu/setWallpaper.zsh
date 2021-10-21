setWallpaper() {
  [[ ! -d ~/Wallpapers ]] && mkdir ~/Wallpapers
  [[ -z "$(ls ~/Wallpapers)" ]] && ${zsb}.throw "You don't have any wallpapers in $(hl ~/Wallpapers)."

  local chosenWallpaper=$(builtin cd ~/Wallpapers && ls | sxiv -ftio | awk 'FNR == 1')

  [[ -z "$chosenWallpaper" ]] && ${zsb}.info "Cancelled." && return 1
  xwallpaper --center ~/Wallpapers/${chosenWallpaper}

  # Persistence after reboot (see .xprofile)
  find ~/temp/wallpaper -lname "${HOME}/Wallpapers/*" -delete
  mkdir -p ~/temp/wallpaper
  ln -s ~/Wallpapers/${chosenWallpaper} ~/temp/wallpaper/${chosenWallpaper}
}
