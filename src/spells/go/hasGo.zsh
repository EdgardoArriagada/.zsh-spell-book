hasGo() {
  printAndRun 'go version' && return 0

  ${zsb}.info "You dont have go, redirecting to `hl "https://go.dev/doc/install"`"
  loading 2
  zsb_open "https://go.dev/doc/install"
}

