setWallpaper() {
  mkdir -pv ~/Wallpapers
  [[ -z "$(ls ~/Wallpapers)" ]] && ${zsb}.throw "You don't have any wallpapers in $(hl ~/Wallpapers)."

  local chosenWallpaper=$(builtin cd ~/Wallpapers && ls | sxiv -ftio | head -1)

  [[ -z "$chosenWallpaper" ]] && ${zsb}.info "Cancelled." && return 1

  xwallpaper --center ~/Wallpapers/${chosenWallpaper} &&
    ${zsb}.success "New wallpaper is $(hl ${chosenWallpaper})"

  # Persistence after reboot (see .xprofile)
  mkdir -pv ~/temp/wallpaper
  find ~/temp/wallpaper -lname "${HOME}/Wallpapers/*" -delete
  ln -s ~/Wallpapers/${chosenWallpaper} ~/temp/wallpaper/
}
