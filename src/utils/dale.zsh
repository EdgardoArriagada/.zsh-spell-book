dale() {
  if ! code --version >/dev/null 2>&1; then
    echo "${ZSB_ERROR} You need to install $(hl vscode) first"
    return 1
  fi

  code . && close
}
