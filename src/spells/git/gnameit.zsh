gnameit() {
  # Gets userName in lowercase
  local userName="${$(git config user.name):l}"
  local oldUrl="$(git config --get remote.origin.url)"

  [[ -z "$userName" ]] &&
    ${zsb}.throw "Username not set."

  [[ -z "$oldUrl" ]] &&
    ${zsb}.throw "Origin url not set."

  [[ "$oldUrl" =~ "^http[s]:\/{2}.+@" ]] &&
    ${zsb}.cancel "Username in url is already set.
    \r$(hl $oldUrl)"

  local regexReplacer="s/:\/{2}/:\/\/${userName}@/"
  local newUrl="$(sed -E "$regexReplacer" <<< "$oldUrl")"

  git remote set-url origin "$newUrl" && \
    ${zsb}.success "Url changed from
      \r${ZSB_COLOR_FOREGROUND_RED} - ${oldUrl}${ZSB_NO_COLOR}
      \r${ZSB_COLOR_FOREGROUND_GREEN} + ${newUrl}${ZSB_NO_COLOR}"
}
