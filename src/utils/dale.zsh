dale() {
  if ! code --version >/dev/null 2>&1; then
    ${zsb}.throw "You need to install `hl vscode` first"
  fi

  code . && (( $ZSB_MACOS )) || close
}
